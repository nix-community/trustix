// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

import type { Component } from "solid-js";
import { A } from "@solidjs/router";

const homepageURL = "https://github.com/nix-community/trustix";
const officialInstanceURL = "https://r13y.trustix.dev";

const About: Component = () => {
  return (
    <>
      <h3 class="text-lg font-bold">FAQ about Trustix R13y</h3>

      <span class="py-4 mx-1">
        <h2 class="text-lg font-bold">What is this?</h2>
        <p>
          This tool aggregates data from multiple <A href={homepageURL}>Trustix</A> log operators (i.e. builders) and
          cross compares them to establish reproducibility.
        </p>
        <p>
          Builds are run on a
          variety of platforms and hardware configurations.
        </p>
      </span>

      <span class="py-4 mx-1">
        <h2 class="text-lg font-bold">
          How confident can we be in the results?
        </h2>
        <p>
          Fairly. We don't know exactly how many reproducibility issues are
          being exercised already, but it's fair to say quite a few. <br />
          The quality of the results depend on the diverseness of builders,
          If all builders run the same CPU model the quality will be lower
        </p>
        <br />
        <p>
          It isn't possible to guarantee a package is reproducible, just like it isn't
          possible to prove software is bug-free. It is possible there is
          nondeterminism in a package source, waiting for some specific
          circumstance.
        </p>
      </span>

      <span class="py-4 mx-1">
        <h2 class="text-lg font-bold">How can I help?</h2>
        <p>
          You can start submitting builds to your own Trustix log and add the
          log to the <A class="underline" href={officialInstanceURL}>official Trustix R13y instance</A>.
          <br />
          It's especially helpful if you are bootstrapping (not using
          cache.nixos.org) for any substitutions, but every bit helps.
        </p>
        <br />
        <p>

          Another thing you can do is to report issues and help to fix them at <A class="underline" href={homepageURL}>{homepageURL}</A>.
        </p>
      </span>
    </>
  );
};

export default About;
