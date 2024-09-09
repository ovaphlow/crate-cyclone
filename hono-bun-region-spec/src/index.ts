import { serve } from '@hono/node-server'
import { Hono } from 'hono'
import * as dotenv from 'dotenv'
import weixinZhCnRouter from './routes/weixin-zhcn'
import emailZhCnRouter from './routes/email-zhcn'

dotenv.config()

const app = new Hono()

app.get('/', (c) => {
    return c.text('Hello Hono!')
})

app.route('/crate-region-spec-api/zhcn/weixin', weixinZhCnRouter)
app.route('/crate-region-spec-api/zhcn/email', emailZhCnRouter)

const port = process.env.REGION_SPEC_HTTP_PORT || 8488
console.log(`Server is running on port ${port}`)

serve({
    fetch: app.fetch,
    port: Number(port),
})
