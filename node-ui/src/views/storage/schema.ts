import { z } from 'zod'

export const pairSchema = z.object({
    Key: z.string(),
    Value: z.string(),
})

export type PairInfo = z.infer<typeof pairSchema>