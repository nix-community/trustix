import { ReproducibilityAPI } from "./api/api_connectweb";
import {
  createConnectTransport,
  createPromiseClient,
} from "@bufbuild/connect-web";

const client = createPromiseClient(
  ReproducibilityAPI,
  createConnectTransport({
    baseUrl: "/api",
  }),
);

export default client;
