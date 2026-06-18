import { spawn } from 'node:child_process'
import { mkdtemp, readFile, rm } from 'node:fs/promises'
import { createWriteStream } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'

const repoRoot = new URL('..', import.meta.url)
const serverDir = new URL('../server/', import.meta.url)
const port = 18080 + Math.floor(Math.random() * 1000)
const addr = `127.0.0.1:${port}`
const baseUrl = `http://${addr}`
const tempDir = await mkdtemp(join(tmpdir(), 'sanjie-smoke-'))
const dbPath = join(tempDir, 'sanjie-smoke.db')
const stdoutPath = join(tempDir, 'server.out.log')
const stderrPath = join(tempDir, 'server.err.log')
const binaryPath = join(tempDir, process.platform === 'win32' ? 'sanjie-smoke.exe' : 'sanjie-smoke')

let stdout = null
let stderr = null
let child = null
let childExit = null

try {
  await runCommand('go', ['build', '-o', binaryPath, '.'], serverDir)
  stdout = createWriteStream(stdoutPath)
  stderr = createWriteStream(stderrPath)
  child = spawn(binaryPath, [], {
    cwd: serverDir,
    env: {
      ...process.env,
      SANJIE_SERVER_ADDR: addr,
      SANJIE_DATABASE_DSN: dbPath,
      SANJIE_DEMO_SEED: 'true'
    },
    stdio: ['ignore', 'pipe', 'pipe'],
    windowsHide: true
  })
  child.stdout.pipe(stdout)
  child.stderr.pipe(stderr)
  child.once('exit', (code, signal) => {
    childExit = { code, signal }
  })

  await waitForHealthy(`${baseUrl}/api/health`, 15000)
  const dashboard = await getJSON(`${baseUrl}/api/dashboard`)
  if (dashboard.code !== 0) {
    throw new Error(`dashboard returned code ${dashboard.code}`)
  }
  const summary = {
    ok: true,
    baseUrl,
    database: dbPath,
    pendingCaptureCount: dashboard.data?.pendingCaptureCount,
    unhandledAlertCount: dashboard.data?.unhandledAlertCount
  }
  console.log(JSON.stringify(summary, null, 2))
} catch (error) {
  await printDiagnostics(error)
  process.exitCode = 1
} finally {
  await stopChild()
  stdout?.close()
  stderr?.close()
  await delay(250)
  await removeTempDir()
}

async function waitForHealthy(url, timeoutMs) {
  const start = Date.now()
  let lastError = null
  while (Date.now() - start < timeoutMs) {
    if (childExit) {
      throw new Error(`server exited early: ${JSON.stringify(childExit)}`)
    }
    try {
      const body = await getJSON(url)
      if (body.code === 0) {
        return
      }
      lastError = new Error(`health returned code ${body.code}`)
    } catch (error) {
      lastError = error
    }
    await delay(500)
  }
  throw new Error(`server did not become healthy: ${lastError?.message || 'timeout'}`)
}

async function getJSON(url) {
  const response = await fetch(url, {
    headers: {
      'X-User-ID': '1'
    }
  })
  if (!response.ok) {
    throw new Error(`${url} returned HTTP ${response.status}`)
  }
  return response.json()
}

async function stopChild() {
  if (!child || childExit) {
    return
  }
  child.kill('SIGTERM')
  for (let i = 0; i < 20; i += 1) {
    if (childExit) {
      return
    }
    await delay(100)
  }
  if (!childExit) {
    child.kill('SIGKILL')
  }
}

async function runCommand(command, args, cwd) {
  await new Promise((resolve, reject) => {
    const proc = spawn(command, args, {
      cwd,
      stdio: ['ignore', 'pipe', 'pipe'],
      windowsHide: true
    })
    let stderrText = ''
    proc.stderr.on('data', chunk => {
      stderrText += chunk.toString()
    })
    proc.once('error', reject)
    proc.once('exit', code => {
      if (code === 0) {
        resolve()
        return
      }
      reject(new Error(`${command} ${args.join(' ')} failed with code ${code}: ${stderrText.trim()}`))
    })
  })
}

async function removeTempDir() {
  let lastError = null
  for (let i = 0; i < 10; i += 1) {
    try {
      await rm(tempDir, { recursive: true, force: true })
      return
    } catch (error) {
      lastError = error
      await delay(250)
    }
  }
  throw lastError
}

async function printDiagnostics(error) {
  const [stdoutText, stderrText] = await Promise.all([
    readFile(stdoutPath, 'utf8').catch(() => ''),
    readFile(stderrPath, 'utf8').catch(() => '')
  ])
  console.error(error)
  console.error('--- server stdout ---')
  console.error(stdoutText.trim())
  console.error('--- server stderr ---')
  console.error(stderrText.trim())
  console.error('--- repo ---')
  console.error(repoRoot.pathname)
}

function delay(ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}
