<template>
    <div class="flex flex-col sm:flex-row items-center justify-center border border-gray-300 rounded-md p-4 min-h-48 sm:min-h-32 gap-4 bg-white">
        <div class="flex flex-col items-center gap-3 text-center w-full sm:flex-1">
            <div class="font-bold text-xl">{{ name }}</div>
            <div>
                <div v-if="isCooked">
                    <span class="font-bold text-blue-600">{{ cookCount }}</span>回食べた
                </div>
                <div v-else>
                    <span class="font-bold text-red-700">まだ食べてない</span>
                </div>
            </div>
            <div v-if="formattedLastCookedDate">
                <span class="font-bold text-blue-500">{{ formattedLastCookedDate }}</span>が最後に食べた日
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

const formatToJstDate = (dateString?: string): string | undefined => {
    if (!dateString) {
        return undefined;
    }
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
        timeZone: 'Asia/Tokyo',
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
    });
};

const isCooked = computed((): boolean => (0 < props.cookCount));
const formattedLastCookedDate = computed((): string | undefined => (
    formatToJstDate(props.lastCookedDate)
));

</script>
