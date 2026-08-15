(() => {
  const PROJECT_URL =
    "https://marinjursic.github.io/production-readiness-checklist/";
  const PROJECT_TEXT =
    "Production Readiness Checklist: 1,421 evidence-driven checks for shipping web applications with confidence.";

  const setStatus = (button, message) => {
    const group = button.closest(".prc-share-group");
    const status = group?.querySelector(".prc-share-status");
    if (!status) return;

    status.textContent = message;
    window.setTimeout(() => {
      if (status.textContent === message) status.textContent = "";
    }, 3000);
  };

  const copyLink = async (url, button) => {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(url);
      setStatus(button, "Link copied");
      return;
    }

    window.prompt("Copy this link:", url);
  };

  const share = async (button) => {
    const isHomepage = document.body.classList.contains("prc-is-homepage");
    const url = isHomepage ? PROJECT_URL : window.location.href;
    const data = {
      title: "Production Readiness Checklist",
      text: PROJECT_TEXT,
      url,
    };

    try {
      if (navigator.share && (!navigator.canShare || navigator.canShare(data))) {
        await navigator.share(data);
        setStatus(button, "Shared");
      } else {
        await copyLink(url, button);
      }
    } catch (error) {
      if (error?.name !== "AbortError") {
        try {
          await copyLink(url, button);
        } catch {
          setStatus(button, "Copy failed");
        }
      }
    }
  };

  const initialize = () => {
    document.body.classList.toggle(
      "prc-is-homepage",
      document.querySelector(".prc-hero") !== null,
    );

    document.querySelectorAll("[data-prc-share]").forEach((button) => {
      button.hidden = false;
      if (button.dataset.prcShareReady === "true") return;
      button.dataset.prcShareReady = "true";
      button.addEventListener("click", () => share(button));
    });
  };

  if (typeof document$ !== "undefined") {
    document$.subscribe(initialize);
  } else if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initialize, { once: true });
  } else {
    initialize();
  }
})();
