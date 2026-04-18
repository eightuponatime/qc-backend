import htmx from "htmx.org"
import { getOrCreateDeviceId, getBrowserInfo } from "./device_data.js"

const deviceId = getOrCreateDeviceId()
const browserInfo = getBrowserInfo()

console.log("device_id:", deviceId)
console.log("browser info:", browserInfo)

document.addEventListener("DOMContentLoaded", () => {
  const container = document.getElementById("today-vote-container")
  if (!container) return

  htmx.ajax("GET", `/fragments/today-vote?device_id=${encodeURIComponent(deviceId)}`, {
    target: "#today-vote-container",
    swap: "innerHTML",
  })
})