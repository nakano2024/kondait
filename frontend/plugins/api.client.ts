import { defineNuxtPlugin, useRuntimeConfig } from '#app';
import { ENV_DEVELOPMENT } from '~/constant/env';
import {
    createRecommendedCookingItemApi,
    createRecommendedCookingItemApiMock,
    type RecommendedCookingItemApi,
} from '~/api/recommended-cooking-item';

export default defineNuxtPlugin(() => {
    const config = useRuntimeConfig();
    const recommendedCookingItemsApi = config.public.env === ENV_DEVELOPMENT
            ? createRecommendedCookingItemApiMock()
            : createRecommendedCookingItemApi();

    return {
        provide: {
            recommendedCookingItemsApi,
        },
    };
});

declare module '#app' {
    interface NuxtApp {
        $recommendedCookingItemsApi: RecommendedCookingItemApi;
    }
}

declare module 'vue' {
    interface ComponentCustomProperties {
        $recommendedCookingItemsApi: RecommendedCookingItemApi;
    }
}
