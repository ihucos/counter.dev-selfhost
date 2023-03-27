customElements.define(
    tagName(),
    class extends HTMLElement {
        connectedCallback() {
            this.innerHTML = `
                    <footer>
                      <section class="footer">
                        <div class="content">
                          <div class="footer-logo">
                            <img
                              src="/img/logotype--gray.svg"
                              width="140"
                              height="32"
                              alt="Logotype"
                            />
                          </div>
                          <nav class="nav-footer">
                          </nav>
                          <div class="footer-contacts">
                            <div class="footer-contacts-social mb16">
                              <a
                                href="https://github.com/ihucos/counter.dev"
                                class="github mr16"
                                target="_blank"
                                rel="nofollow"
                              ></a>
                              <a href="https://twitter.com/DevCounter" class="twitter" target="_blank" rel="nofollow"></a>
                            </div>
                            <div class="caption gray">
                              You are logged in as asdf – <a href="/logout" class="caption gray underline">Logout</a>
                            </div>
                          </div>
                        </div>
                      </section>

                    </footer>`;
        }
    }
);
