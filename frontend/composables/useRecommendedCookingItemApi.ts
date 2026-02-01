import { useNuxtApp } from '#app';
import type { RecommendedCookingItemApi } from '~/api/recommended-cooking-item';

export const useRecommendedCookingItemApi = (): RecommendedCookingItemApi => {
    const { $recommendedCookingItemsApi } = useNuxtApp();

    return $recommendedCookingItemsApi;
};
