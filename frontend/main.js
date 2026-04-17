import "htmx.org"
import { getDeviceData } from "./device_data.js"

(async () => {
    const data = await getDeviceData()
    console.log(data)
})()