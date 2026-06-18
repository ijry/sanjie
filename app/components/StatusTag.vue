<template>
  <up-tag :text="labelOf(value)" :type="typeOf(value)" size="mini" />
</template>

<script setup>
import { labelOf, typeOf } from '../utils/status'

defineProps({
  value: {
    type: String,
    default: ''
  }
})
</script>

