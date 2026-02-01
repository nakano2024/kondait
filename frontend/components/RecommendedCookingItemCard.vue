<template>
    <div class="flex flex-col sm:flex-row items-center justify-center border border-gray-300 rounded-md p-4 min-h-32 gap-4 bg-white">
        <div class="flex flex-col items-center gap-3 text-center w-full sm:flex-1">
            <div class="font-bold text-xl">{{ name }}</div>
            <div>
                <div v-if="isCooked">
                    <span class="font-bold text-blue-600">{{ cookCount }}</span>回食べた
                </div>
                <div v-else>
                    <span class="font-bold text-red-700">一回も食べてない</span>
                </div>
            </div>
            <div v-if="lastCookedDate">
                <span class="font-bold text-blue-500">{{ lastCookedDate }}</span>が最後に食べた日
            </div>
        </div>
        <div class="flex items-center justify-center">
            <CookButton />
        </div>
    </div>
</template>

<script setup lang="ts">

const props = defineProps<{
    code: string;
    name: string;
    cookCount: number;
    lastCookedDate?: string;
}>();

const isCooked = computed((): boolean => (0 < props.cookCount));

</script>
