// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

import type { Component } from "solid-js";
import { lazy, onMount, For } from "solid-js";
import { Routes, Route, A } from "@solidjs/router";
import { themeChange } from "theme-change";
import { OcGitbranch2 } from "solid-icons/oc";

const Derivation = lazy(() => import("./pages/drv"));
const Attrs = lazy(() => import("./pages/attrs"));
const Diff = lazy(() => import("./pages/diff"));
const About = lazy(() => import("./about"));

const gitRepoURL = "https://github.com/nix-community/trustix";

const links: {
  title: string;
  href: string;
}[] = [
  {
    title: "Attributes",
    href: "/attrs",
  },
];

const App: Component = () => {
  onMount(async () => {
    themeChange();
  });

  const aboutModalInput = (
    <input type="checkbox" id="about-modal" class="modal-toggle" />
  );

  return (
    <>
      {aboutModalInput}
      <label for="about-modal" class="modal cursor-pointer">
        <label class="modal-box relative" for="">
          <About />
        </label>
      </label>

      <div class="bg-base-200 min-h-screen">
        <nav id="main-nav" class="bg-primary p-2">
          <A href="/">
            <span>Trustix r13y</span>
          </A>

          <For each={links}>
            {(l) => (
              <A class="text-primary-content m-2" href={l.href}>
                {l.title}
              </A>
            )}
          </For>

          <div class="float-right">
            <span>
              <A href={gitRepoURL}>
                <button>
                  <OcGitbranch2 color="white" />
                </button>
              </A>
            </span>

            <button data-toggle-theme="dark,light" data-act-class="ACTIVECLASS">
              <span class="" data-tip="toggle dark mode">
                💡
              </span>
            </button>

            <span>
              <button
                onClick={() => {
                  aboutModalInput.checked = true;
                }}
              >
                ❓
              </button>
            </span>
          </div>
        </nav>

        <div class="flex justify-evenly place-items-center">
          <Routes>
            <Route
              path="/"
              component={() => {
                return (
                  <div class="w-3/4">
                    <About />
                  </div>
                );
              }}
            />
            <Route path="/attrs" component={Attrs} />
            <Route path="/drv" component={Derivation} />
            <Route path="/diff" component={Diff} />
          </Routes>
        </div>
      </div>
    </>
  );
};

export default App;
