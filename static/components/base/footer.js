customElements.define(
    tagName(),
    class extends HTMLElement {
        connectedCallback() {
            this.innerHTML = `
                    <footer>
                      <section class="footer">
                        <div class="content">
                          <div class="footer-logo">
                            <a href="https://github.com/ihucos/counter.dev-selfhosted" target="_blank">
                                <img
                                  src="/img/logotype--gray.svg"
                                  width="140"
                                  height="32"
                                  alt="Logotype"
                                />
                            </a>
                            <div class="caption gray">
                                Leave a <a href="#">tip</a>!
                            </div>
                          </div>
                        </div>
                      </section>

                    </footer>`;
        }
    }
);
