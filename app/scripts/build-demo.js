#!/usr/bin/env node

const { execSync } = require('child_process')
const fs = require('fs')
const path = require('path')

const appDir = path.resolve(__dirname, '..')
const repoRoot = path.resolve(appDir, '..')
const manifestPath = path.join(appDir, 'manifest.json')
const outputDir = process.argv[2] || path.join(repoRoot, 'docs-site', 'docs', 'public', 'h5-demo')

console.log('[sanjie-demo] Building embedded H5 mock demo')
console.log('[sanjie-demo] Output:', outputDir)

const originalManifest = fs.readFileSync(manifestPath, 'utf8')
const manifest = JSON.parse(originalManifest)

if (!manifest.h5) manifest.h5 = {}
if (!manifest.h5.router) manifest.h5.router = {}
manifest.h5.router.mode = 'hash'
manifest.h5.router.base = './'
manifest.h5.publicPath = './'

fs.writeFileSync(manifestPath, `${JSON.stringify(manifest, null, 2)}\n`, 'utf8')

try {
  const isWindows = process.platform === 'win32'
  const uniBin = path.join(appDir, 'node_modules', '.bin', isWindows ? 'uni.cmd' : 'uni')
  execSync(`"${uniBin}" build`, {
    cwd: appDir,
    stdio: 'inherit',
    env: {
      ...process.env,
      UNI_INPUT_DIR: appDir,
      VITE_SANJIE_MOCK: 'true'
    }
  })

  const buildDir = path.join(appDir, 'dist', 'build', 'h5')
  if (!fs.existsSync(buildDir)) {
    throw new Error(`H5 build output not found: ${buildDir}`)
  }

  fs.rmSync(outputDir, { recursive: true, force: true })
  fs.mkdirSync(outputDir, { recursive: true })
  copyDir(buildDir, outputDir)
  console.log('[sanjie-demo] Embedded H5 mock demo ready')
} finally {
  fs.writeFileSync(manifestPath, originalManifest, 'utf8')
  console.log('[sanjie-demo] Restored manifest.json')
}

function copyDir(src, dest) {
  for (const entry of fs.readdirSync(src, { withFileTypes: true })) {
    const srcPath = path.join(src, entry.name)
    const destPath = path.join(dest, entry.name)
    if (entry.isDirectory()) {
      fs.mkdirSync(destPath, { recursive: true })
      copyDir(srcPath, destPath)
    } else {
      fs.copyFileSync(srcPath, destPath)
    }
  }
}
