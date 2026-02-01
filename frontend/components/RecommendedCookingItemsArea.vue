<template>
    <div>
        <div v-if="isLoading">
            <RecommendedCookingItemSkeleton
                v-for="n in 5"
                :key="n"
            />
        </div>
        <div v-else>
            <div v-if="cookingItems.length" class="recommended-cooking-items">
                <RecommendedCookingItemCard
                    v-for="item in cookingItems"
                    :key="item.code"
                    :code="item.code"
                    :name="item.name"
                    :cook-count="item.cookCount"
                    :last-cooked-date="item.formattedLastCookedDate"
                />
            </div>
            <div v-else>
                おすすめ可能な献立がありません。
                献立を追加
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRecommendedCookingItemApi } from '~/composables/useRecommendedCookingItemApi';

interface CookingItem {
    code: string,
    name: string,
    cookCount: number,
    formattedLastCookedDate?: string,
}

const isLoading = ref<boolean>(true);
const cookingItems = ref<CookingItem[]>([]);
const recommendedCookingItemApi = useRecommendedCookingItemApi();

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

function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}

const fetchCookingItems = async () => {
    isLoading.value = true;
    const recommendedCookingItems = await recommendedCookingItemApi.getRecommendedCookingItems();
    cookingItems.value = recommendedCookingItems.map((item) => ({
        code: item.code,
        name: item.name,
        cookCount: item.cookCount,
        formattedLastCookedDate: formatToJstDate(item.lastCookedDate),
    }));
    await sleep(3000);
    isLoading.value = false;
};

onMounted(() => {
    fetchCookingItems();
});
</script>

<style scoped>
</style>
