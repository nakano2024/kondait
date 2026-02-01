import type { RecommendedCookingItem } from "../type/cooking-item";

export interface RecommendedCookingItemApi {
    getRecommendedCookingItems: () => Promise<RecommendedCookingItem[]>
}

const sleep = (ms: number): Promise<void> => {
  return new Promise(resolve => setTimeout(resolve, ms));
}

export const createRecommendedCookingItemApiMock = (): RecommendedCookingItemApi => {
    return {
        getRecommendedCookingItems: async (): Promise<RecommendedCookingItem[]> => {
            await sleep(2000);
            return [
                {
                    code: '5b3b0b2a-1a9b-4df0-9f6b-0c1e6c3b8f01',
                    name: '納豆ご飯',
                    cookCount: 0,
                    lastCookedDate: undefined,
                },
                {
                    code: 'b1c9c9b2-4b8b-4f9c-8d7d-38d1f7e1a2f2',
                    name: '味噌汁',
                    cookCount: 3,
                    lastCookedDate: '2024-01-18T00:00:00Z',
                },
                {
                    code: '2f8b7d73-1c6a-4e4a-9d53-2c2d1b0c6b03',
                    name: '鶏むね肉の照り焼き',
                    cookCount: 1,
                    lastCookedDate: '2024-01-12T00:00:00Z',
                },
                {
                    code: '9e1a4f2b-7d3b-4b7a-8c4e-1f2d3a4b5c04',
                    name: '野菜サラダ',
                    cookCount: 5,
                    lastCookedDate: '2024-01-10T00:00:00Z',
                },
                {
                    code: 'd6f3a2c1-7f5b-4c9b-9a2b-6b7c8d9e0f05',
                    name: 'カレー',
                    cookCount: 2,
                    lastCookedDate: '2024-01-08T00:00:00Z',
                },
            ]
        }
    }
}

export const createRecommendedCookingItemApi = (): RecommendedCookingItemApi => {
    return {
        getRecommendedCookingItems: async (): Promise<RecommendedCookingItem[]> => {
            try {
                // APIリクエストロジックを実装
                return []
            }
            catch (error) {
                const message = error instanceof Error ? error.message : String(error);
                console.error('データの取得に失敗しました。', error);
                throw new Error(`データの取得に失敗しました。 ${message}`);
            }
        }
    }
}
