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
    startMealCountdowns()
    initMealVoteForms()
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

let mealCountdownInterval = null

function startMealCountdowns() {
  if (mealCountdownInterval) {
    clearInterval(mealCountdownInterval)
    mealCountdownInterval = null
  }

  const container = document.getElementById("vote-ui-container")
  if (!container) return

  const cards = container.querySelectorAll("[data-minutes-until-open][data-server-now]")
  if (!cards.length) return

  let hasReloaded = false

  const updateAllCountdowns = () => {
    for (const card of cards) {
      const minutesUntilOpen = Number(card.getAttribute("data-minutes-until-open"))
      const serverNowIso = card.getAttribute("data-server-now")
      const countdownEl = card.querySelector(".meal-card__countdown")
      const statusEl = card.querySelector(".meal-card__status--countdown")

      if (!countdownEl || Number.isNaN(minutesUntilOpen) || !serverNowIso) continue

      const serverNow = new Date(serverNowIso)
      const targetTime = new Date(serverNow.getTime() + minutesUntilOpen * 60 * 1000)
      const now = new Date()
      const diffMs = targetTime.getTime() - now.getTime()

      if (diffMs <= 0) {
        countdownEl.textContent = t("ui.available_now", "Уже доступно")

        if (!hasReloaded) {
          hasReloaded = true
          clearInterval(mealCountdownInterval)
          mealCountdownInterval = null

          loadVoteUI().then(() => {
            startMealCountdowns()
          })
        }

        continue
      }

      const totalMinutes = Math.ceil(diffMs / 60000)
      const hours = Math.floor(totalMinutes / 60)
      const minutes = totalMinutes % 60

      if (hours > 0) {
        countdownEl.textContent = t(
          "ui.countdown_hours_minutes",
          "Осталось: {hours} ч {minutes} мин"
        )
          .replace("{hours}", String(hours))
          .replace("{minutes}", String(minutes))
      } else {
        countdownEl.textContent = t(
          "ui.countdown_minutes",
          "Осталось: {minutes} мин"
        ).replace("{minutes}", String(minutes))
      }

      if (totalMinutes <= 10) {
        card.classList.add("meal-card--soon")
        if (statusEl) {
          statusEl.textContent = t("ui.opens_soon", "Скоро откроется")
        }
      } else {
        card.classList.remove("meal-card--soon")
        if (statusEl) {
          statusEl.textContent = t("ui.not_available_yet", "Пока недоступно")
        }
      }
    }
  }

  updateAllCountdowns()
  mealCountdownInterval = setInterval(updateAllCountdowns, 60000)
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
    startMealCountdowns()
    initMealVoteForms()
  }
})

document.body.addEventListener("change", (event) => {
  if (!event.target.matches(".js-meal-vote-form input[name='rating']")) return

  const form = event.target.closest(".js-meal-vote-form")
  if (!form) return

  updateMealVoteButtonState(form)
})

document.body.addEventListener("input", (event) => {
  if (!event.target.matches(".js-meal-vote-form textarea[name='review']")) return

  const form = event.target.closest(".js-meal-vote-form")
  if (!form) return

  updateMealVoteButtonState(form)
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

function initMealVoteForms() {
  const forms = document.querySelectorAll(".js-meal-vote-form")
  for (const form of forms) {
    updateMealVoteButtonState(form)
  }
}

function updateMealVoteButtonState(form) {
  const submitButton = form.querySelector(".meal-card__button[type='submit']")
  if (!submitButton) return

  const selectedRating = form.querySelector("input[name='rating']:checked")?.value || ""
  const review = form.querySelector("textarea[name='review']")?.value || ""
  const hasExistingVote = form.dataset.hasExistingVote === "true"

  if (!selectedRating) {
    submitButton.disabled = true
    return
  }

  if (!hasExistingVote) {
    submitButton.disabled = false
    return
  }

  const initialRating = form.dataset.initialRating || ""
  const initialReview = form.dataset.initialReview || ""
  submitButton.disabled = selectedRating === initialRating && review === initialReview
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
