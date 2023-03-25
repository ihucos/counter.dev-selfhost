customElements.define(
    tagName(),
    class extends HTMLElement {
        connectedCallback() {
            this.innerHTML = `
              <a href="/logout" class="btn-secondary btn-icon">
                <img src="/img/logout.svg" width="24" height="24" alt="Settings" />
                </a>`;
        }
    }
);
