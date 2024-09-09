import { serve } from '@hono/node-server'
import axios from 'axios'
import { Hono } from 'hono'
import * as dotenv from 'dotenv'
import weixinZhCnRouter from './routes/weixin-zhcn'
import emailZhCnRouter from './routes/email-zhcn'

dotenv.config()

const initHQ = async () => {
    const addr = process.env.HQ_ADDR || '127.0.0.1'
    const port = process.env.HQ_PORT || 8421
    const body = {
        name: 'crate-region-spec-api',
        protocol: 'http',
        host: '127.0.0.1',
        port: 8488,
        healthCheck: {
            endpoint: '/healthCheck'
        }
    }
    try {
        const response = await axios.post(`http://${addr}:${port}/crate-hq-api/service`, body)
        if (response.status < 400) {
            console.info('initHQ success')
        }
    } catch (error) {
        console.error(error)
    }
}
initHQ()

const app = new Hono()

app.use('*', async (c, next) => {
    const { method, url } = c.req
    const start = Date.now()
    await next()
    const ms = Date.now() - start
    console.log(`${method} ${url} - ${ms}ms`)
})

app.get('/', (c) => {
    return c.text('Hello Hono!')
})

app.get('/healthCheck', (c) => {
    return c.text('OK')
})

app.route('/crate-region-spec-api/zhcn/weixin', weixinZhCnRouter)
app.route('/crate-region-spec-api/zhcn/email', emailZhCnRouter)

const port = process.env.REGION_SPEC_HTTP_PORT || 8488
console.log(`Server is running on port ${port}`)

serve({
    fetch: app.fetch,
    port: Number(port),
})
