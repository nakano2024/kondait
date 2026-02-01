<template>
    <h1 class="text-center font-bold text-2xl text-gray-700 mb-5">{{ text }}</h1>
</template>

<script setup lang="ts">
defineProps<{
    text: string;
}>();
</script>
