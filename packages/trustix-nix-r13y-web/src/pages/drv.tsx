import {
  lazy,
  Component,
  createSignal,
  createResource,
  For,
  Show,
  Suspense,
} from "solid-js";
import { Routes, Route, useParams, useSearchParams, A } from "@solidjs/router";

import {
  createConnectTransport,
  createPromiseClient,
} from "@bufbuild/connect-web";

import { ReproducibilityAPI } from "../api/api_connectweb";
import {
  DerivationReproducibilityRequest,
  DerivationReproducibilityResponse,
  DerivationReproducibilityResponse_Derivation,
} from "../api/api_pb";
import { NameValuePair } from "../lib";

type DerivationReproducibilityPaths = {
  [key: string]: DerivationReproducibilityResponse_Derivation;
};

const loading = <h1>Loading...</h1>;

const fetchDerivationReproducibility = async (
  drvPath
): DerivationReproducibilityResponse => {
  const client = createPromiseClient(
    ReproducibilityAPI,
    createConnectTransport({
      baseUrl: "/api",
    })
  );

  const req = new DerivationReproducibilityRequest({
    DrvPath: drvPath,
  });

  return await client.derivationReproducibility(req);
};

const renderPaths = (
  title: string,
  paths: DerivationReproducibilityPaths
): Component => {
  if (Object.keys(paths).length == 0) {
    return <></>;
  }

  const derivations = NameValuePair.fromMap(paths);

  return (
    <>
      <div className="divider"></div>

      <div className="grid place-items-center w-11/12">
        <h3 className="text-xl font-bold underline">{title}</h3>

        <ul className="grid auto-cols-auto list-inside">
          <For each={derivations}>
            {({ name, value }) => {
              const drvPath = name;
              const drvOutputs = NameValuePair.fromMap(value.Outputs);

              return (
                <li>
                  <div className="card drv-card bg-base-100 shadow-xl m-2">
                    <div className="card-body">
                      <A
                        href={`/drv?storePath=${encodeURIComponent(
                          encodeURIComponent(drvPath)
                        )}`}
                      >
                        <h2
                          className="card-title tooltip"
                          data-tip="Derivation store path"
                        >
                          {drvPath}
                        </h2>
                      </A>

                      <ul class="list-disc list-inside">
                        <For each={drvOutputs}>
                          {({ name, value }) => {
                            const output = name;
                            const storePath = value.StorePath;

                            // Sort output hashes according to popularity
                            const outputHashes = NameValuePair.fromMap(
                              value.OutputHashes
                            ).sort(
                              (a, b) =>
                                a.value.LogIDs.length > b.value.LogIDs.length
                            );

                            return (
                              <li>
                                <span>
                                  <span
                                    className="font-bold tooltip"
                                    data-tip="Output name"
                                  >
                                    {output}
                                  </span>
                                  &nbsp;
                                  <span
                                    className="tooltip"
                                    data-tip="Output store path"
                                  >
                                    ({storePath})
                                  </span>
                                </span>

                                <ul class="list-disc list-inside">
                                  <For each={outputHashes}>
                                    {({ name, value }) => {
                                      const outputNarinfoHash = name;
                                      const logIDs: Array<Number> =
                                        value.LogIDs.map(Number);

                                      return (
                                        <li>
                                          <p
                                            className="tooltip"
                                            data-tip="Narinfo hash"
                                          >
                                            {outputNarinfoHash}: {logIDs}
                                          </p>
                                        </li>
                                      );
                                    }}
                                  </For>
                                </ul>
                              </li>
                            );
                          }}
                        </For>
                      </ul>
                    </div>
                  </div>
                </li>
              );
            }}
          </For>
        </ul>
      </div>
    </>
  );
};

const Derivation: Component = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [drvReprod] = createResource(
    () => searchParams.storePath,
    fetchDerivationReproducibility
  );

  return (
    <>
      <div>
        <h2
          className="tooltip text-xl font-bold"
          data-tip="Derivation store path"
        >
          {searchParams.storePath}
        </h2>

        <Suspense fallback={loading}>
          <Show when={drvReprod()}>
            {renderPaths("Unreproduced paths", drvReprod()?.UnreproducedPaths)}
          </Show>

          <Show when={drvReprod()}>
            {renderPaths("Reproduced paths", drvReprod()?.ReproducedPaths)}
          </Show>

          <Show when={drvReprod()}>
            {renderPaths(
              "Unknown paths (only built by one log)",
              drvReprod()?.UnknownPaths
            )}
          </Show>

          <Show when={drvReprod()}>
            {renderPaths(
              "Missing paths (not built by any known log)",
              drvReprod()?.MissingPaths
            )}
          </Show>
        </Suspense>
      </div>
    </>
  );
};

export default Derivation;
