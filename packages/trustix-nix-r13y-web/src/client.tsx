// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

import { ReproducibilityAPI } from "./api/api_connectweb";
import {
  createPromiseClient,
} from "@connectrpc/connect";
import {
  createConnectTransport,
} from "@connectrpc/connect-web";

const client = createPromiseClient(
  ReproducibilityAPI,
  createConnectTransport({
    baseUrl: "/api",
  }),
);

export default client;
