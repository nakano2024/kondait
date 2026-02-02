<template>
    <div>
        <SectionHeading text="本日のおすすめ" />
        <div class="content-scroll rounded-md" :class="overflowStyle">
            <RecommendedCookingItemSkeleton
                v-if="isLoading"
                v-for="n in 5"
                :key="n"
            />
            <RecommendedCookingItemCard
                v-else-if="cookingItems.length"
                v-for="item in cookingItems"
                :key="item.code"
                :code="item.code"
                :name="item.name"
                :cook-count="item.cookCount"
                :last-cooked-date="item.lastCookedDate"
            />
            <RecommendedCookingItemsEmpty v-else />
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRecommendedCookingItemApi } from '~/composables/useRecommendedCookingItemApi';
import type { RecommendedCookingItem } from '~/type/cooking-item';

const isLoading = ref<boolean>(true);
const cookingItems = ref<RecommendedCookingItem[]>([]);
const recommendedCookingItemApi = useRecommendedCookingItemApi();

const fetchCookingItems = async () => {
    isLoading.value = true;
    const recommendedCookingItems = await recommendedCookingItemApi.getRecommendedCookingItems();
    cookingItems.value = recommendedCookingItems;
    isLoading.value = false;
};

const overflowStyle = computed((): string => {
    if (isLoading.value) {
        return 'overflow-y-hidden';
    }
    return 'overflow-y-auto';
});

onMounted(() => {
    fetchCookingItems();
});
</script>

<style scoped>
.content-scroll {
    height: 80vh;
    scrollbar-width: none;
    -ms-overflow-style: none;
}

.content-scroll::-webkit-scrollbar {
    display: none;
}
</style>
