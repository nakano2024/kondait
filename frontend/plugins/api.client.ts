import { defineNuxtPlugin, useRuntimeConfig } from '#app';
import { ENV_DEVELOPMENT } from '~/constant/env';
import {
    createRecommendedCookingItemsApi,
    createRecommendedCookingItemsApiMock,
    type RecommendedCookingItemsApi,
} from '~/api/recommended-cooking-item';

export default defineNuxtPlugin(() => {
    const config = useRuntimeConfig();
    const recommendedCookingItemsApi =
        config.public.env === ENV_DEVELOPMENT
            ? createRecommendedCookingItemsApiMock()
            : createRecommendedCookingItemsApi();

    return {
        provide: {
            recommendedCookingItemsApi,
        },
    };
});

declare module '#app' {
    interface NuxtApp {
        $recommendedCookingItemsApi: RecommendedCookingItemsApi;
    }
}

declare module 'vue' {
    interface ComponentCustomProperties {
        $recommendedCookingItemsApi: RecommendedCookingItemsApi;
    }
}
