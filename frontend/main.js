import "./styles/index.css"
import htmx from "htmx.org"
import { getOrCreateDeviceId, getBrowserInfo } from "./device_data.js"

const deviceId = getOrCreateDeviceId()
const browserInfo = getBrowserInfo()

function t(key, fallback = "") {
  const dict = window.__i18n || {}
  return dict[key] || fallback || key
}

window.__deviceInfo = {
  deviceId,
  phoneModel: browserInfo.platform || "Unknown",
  browser: detectBrowser(browserInfo.userAgent),
}

document.addEventListener("DOMContentLoaded", () => {
  const container = document.getElementById("vote-ui-container")
  if (!container) return

  loadVoteUI().then(() => {
    initVoteForm()
  })
})

function loadVoteUI() {
  const params = new URLSearchParams({
    device_id: window.__deviceInfo.deviceId,
    phone_model: window.__deviceInfo.phoneModel,
    browser: window.__deviceInfo.browser,
  })

  return htmx.ajax("GET", `/fragments/vote-ui?${params.toString()}`, {
    target: "#vote-ui-container",
    swap: "innerHTML",
  })
}

function detectBrowser(userAgent) {
  if (userAgent.includes("Edg")) return "Edge"
  if (userAgent.includes("OPR")) return "Opera"
  if (userAgent.includes("Chrome")) return "Chrome"
  if (userAgent.includes("Firefox")) return "Firefox"
  if (userAgent.includes("Safari")) return "Safari"
  return "Unknown"
}

document.body.addEventListener("htmx:sendError", () => {
  showGlobalError(t("errors.network_unreachable", "Нет соединения с сервером."))
})

document.body.addEventListener("htmx:responseError", (event) => {
  const xhr = event.detail.xhr

  if (!xhr) {
    showGlobalError(t("errors.unknown_response_error", "Ошибка ответа сервера."))
    return
  }

  if (xhr.status >= 500) {
    showGlobalError(t("errors.server_unavailable", "На сервере произошла ошибка."))
    return
  }

  if (xhr.status >= 400) {
    showGlobalError(t("errors.request_failed", "Не удалось выполнить запрос."))
    return
  }

  showGlobalError(t("errors.unexpected_client_error", "Произошла непредвиденная ошибка."))
})

document.body.addEventListener("htmx:afterSwap", (event) => {
  if (event.target && event.target.id === "vote-ui-container") {
    hideGlobalError()
    initVoteForm()
  }
})

document.body.addEventListener("change", (event) => {
  if (!event.target.matches(".js-vote-form input[type='radio']")) return

  const form = event.target.closest(".js-vote-form")
  if (!form) return

  updateVoteFormButtonState(form)
})

document.body.addEventListener("input", (event) => {
  if (!event.target.matches(".js-vote-form textarea")) return

  const form = event.target.closest(".js-vote-form")
  if (!form) return

  updateVoteFormButtonState(form)
})

function showGlobalError(message) {
  const container = document.getElementById("global-error")
  if (!container) return

  container.textContent = message
  container.hidden = false
}

function hideGlobalError() {
  const container = document.getElementById("global-error")
  if (!container) return

  container.hidden = true
  container.textContent = ""
}

function initVoteForm() {
  const form = document.querySelector(".js-vote-form")
  if (!form) return

  updateVoteFormButtonState(form)
}

function updateVoteFormButtonState(form) {
  const submitButton = form.querySelector(".vote-submit-bar__button[type='submit']")
  if (!submitButton) return
  const pendingLabel = form.querySelector(".vote-submit-bar__pending")
  const successLabel = form.querySelector(".vote-submit-bar__success")

  const sections = form.querySelectorAll(".js-meal-vote-section")
  let hasChanges = false
  let hasIncompleteSection = false

  for (const section of sections) {
    const mealType = section.dataset.mealType
    const initialRating = section.dataset.initialRating || ""
    const initialReview = section.dataset.initialReview || ""
    const selectedRating = form.querySelector(`input[name='${mealType}_rating']:checked`)?.value || ""
    const review = form.querySelector(`textarea[name='${mealType}_review']`)?.value || ""
    const trimmedReview = review.trim()
    const hasAnyContent = selectedRating !== "" || trimmedReview !== ""
    const isChanged = selectedRating !== initialRating || review !== initialReview

    if (isChanged && !selectedRating && trimmedReview !== "") {
      hasIncompleteSection = true
    }

    if (isChanged && hasAnyContent && selectedRating) {
      hasChanges = true
    }
  }

  submitButton.disabled = !hasChanges || hasIncompleteSection

  if (pendingLabel) {
    pendingLabel.hidden = !(hasChanges || hasIncompleteSection)
  }

  if (successLabel && (hasChanges || hasIncompleteSection)) {
    successLabel.hidden = true
  }
}

function dismissAccessWarning(button) {
  const warningBlock = button?.closest(".access-warning")
  if (!warningBlock) return

  const mealCard = warningBlock.closest(".meal-card")
  if (!mealCard) return

  warningBlock.remove()

  const form = mealCard.querySelector("form")
  if (form) {
    form.hidden = false
  }
}

window.dismissAccessWarning = dismissAccessWarning
