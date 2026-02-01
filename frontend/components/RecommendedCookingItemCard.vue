<template>
    <div class="flex flex-col sm:flex-row items-center justify-center border border-gray-300 rounded-md p-4 min-h-32 gap-4">
        <div class="flex flex-col items-center gap-3 text-center w-full sm:flex-1">
            <div class="font-bold text-xl">{{ name }}</div>
            <div>
                <span class="font-bold text-red-600">{{ cookCount }}</span>回作った
            </div>
            <div v-if="lastCookedDate">
                <span class="font-bold text-sky-600">{{ lastCookedDate }}</span>が最後に作った日
            </div>
        </div>
        <div class="flex items-center justify-center">
            <CookButton />
        </div>
    </div>
</template>

<script setup lang="ts">

defineProps<{
    code: string;
    name: string;
    cookCount: number;
    lastCookedDate?: string;
}>();
</script>
