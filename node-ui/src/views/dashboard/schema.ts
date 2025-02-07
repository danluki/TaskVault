import { z } from 'zod'

const tagSchema = z.object({
    dc: z.string(),
    expect: z.string(),
    port: z.number(),
    region: z.string(),
    role: z.string(),
    rpc_addr: z.string(),
    server: z.string(),
    version: z.string(),
})

export const taskSchema = z.object({
    Name: z.string(),
    Addr: z.string(),
    Port: z.string(),
    Tags: tagSchema,
    statusText: z.string(),
})

export type NodeInfo = z.infer<typeof taskSchema>