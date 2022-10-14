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
  DerivationReproducibilityResponse_DerivationOutput,
  DerivationReproducibilityResponse_DerivationOutputHash,
} from "../api/api_pb";
import { NameValuePair } from "../lib";

type DerivationReproducibilityPaths = {
  [key: string]: DerivationReproducibilityResponse_Derivation;
};
type Logs = { [key: string]: Log };

const loading = <h1>Loading...</h1>;

const fetchDerivationReproducibility = async (
  drvPath,
): DerivationReproducibilityResponse => {
  const client = createPromiseClient(
    ReproducibilityAPI,
    createConnectTransport({
      baseUrl: "/api",
    }),
  );

  const req = new DerivationReproducibilityRequest({
    DrvPath: drvPath,
  });

  return await client.derivationReproducibility(req);
};

const renderDerivationOutput = (
  output: string,
  storePath: string,
  outputHashes: NameValuePair<DerivationReproducibilityResponse_DerivationOutputHash>[],
  logs: Logs,
): Component => {
  // Keeps track of which checkboxes (ie active diffs) are checked
  const checkedNarinfoHashes = new Set<string>();

  const renderOutputHash = (
    outputNarinfoHash: string,
    logIDs: Array<string>,
  ): Component => {
    // Args passed to the input checkbox
    const checkboxArgs = {};

    // Hide the compare checkbox if there is one or less output hash,(nohing to compare)
    if (logIDs.length < 2) {
      checkboxArgs["disabled"] = "disabled";
    }

    const onChecked = (e) => {
      if (e.target.checked) {
        checkedNarinfoHashes.add(outputNarinfoHash);
      } else {
        checkedNarinfoHashes.delete(outputNarinfoHash);
      }
    };

    return (
      <>
        <tr>
          <th>
            <label>
              <input
                onInput={onChecked}
                type="checkbox"
                className="checkbox"
                {...checkboxArgs}
              />
            </label>
          </th>
          <td>
            <div className="flex items-center space-x-3">
              <div>
                <div className="text-sm opacity-50">{outputNarinfoHash}</div>
              </div>
            </div>
          </td>
          <td>
            <For each={logIDs}>
              {(logID) => (
                <>
                  <span className="badge badge-ghost badge-sm">
                    {logs[logID].Name}
                  </span>
                  <br />
                </>
              )}
            </For>
          </td>
        </tr>
      </>
    );
  };

  const onNarinfoClicked = (e) => {
    const checked = checkedNarinfoHashes;

    if (checked.size == 0) {
      alert("No Narinfo hashes selected");
      return;
    } else if (checked.size == 1 || checked.size > 2) {
      alert(
        "Invalid number of Narinfo hashes selected, we can only compare 2 at a time",
      );
      return;
    }

    const [a, b] = checked;

    alert("TODO: Redirect to diff view");
  };

  return (
    <>
      <div className="card bg-base-200 shadow-xl m-3">
        <div className="card-body">
          <h2 className="card-title tooltip" data-tip="Output name">
            {output}
          </h2>
          <p className="font-bold">{storePath}</p>

          {outputHashes.length > 0 && (
            <>
              <div className="overflow-x-auto w-full">
                <table className="table w-full">
                  <thead>
                    <tr>
                      <th>✓</th>
                      <th>Narinfo hash</th>
                      <th>Logs</th>
                    </tr>
                  </thead>
                  <tbody>
                    <For each={outputHashes}>
                      {({ name, value }) =>
                        renderOutputHash(name, value.LogIDs)
                      }
                    </For>
                  </tbody>
                </table>

                {/* show the compare button if there are more than one output hash for the same output */}
                {outputHashes.length > 1 && (
                  <button
                    onClick={onNarinfoClicked}
                    className="btn btn-info btn-sm"
                  >
                    Compare outputs
                  </button>
                )}
              </div>
            </>
          )}
        </div>
      </div>
    </>
  );
};

const renderDerivationOutputs = (
  drvPath: string,
  cardBackground: string,
  drvOutputs: NameValuePair<DerivationReproducibilityResponse_DerivationOutput>[],
  logs: Logs,
): Component => {
  return (
    <>
      <div className={`card drv-card shadow-xl m-2 ${cardBackground}`}>
        <div className="card-body">
          <A
            href={`/drv?storePath=${encodeURIComponent(
              encodeURIComponent(drvPath),
            )}`}
          >
            <h2 className="card-title tooltip" data-tip="Derivation store path">
              {drvPath}
            </h2>
          </A>

          <div>
            <For each={drvOutputs}>
              {({ name, value }) =>
                renderDerivationOutput(
                  name,
                  value.StorePath,
                  NameValuePair.fromMap(value.OutputHashes).sort(
                    (a, b) => a.value.LogIDs.length > b.value.LogIDs.length,
                  ),
                  logs,
                )
              }
            </For>
          </div>
        </div>
      </div>
    </>
  );
};

const renderPaths = (
  title: string,
  cardBackground: string,
  paths: DerivationReproducibilityPaths,
  logs: Logs,
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

        <For each={derivations}>
          {({ name, value }) =>
            renderDerivationOutputs(
              name,
              cardBackground,
              NameValuePair.fromMap(value.Outputs),
              logs,
            )
          }
        </For>
      </div>
    </>
  );
};

const Derivation: Component = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [drvReprod] = createResource(
    () => searchParams.storePath,
    fetchDerivationReproducibility,
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
            {renderPaths(
              "Unreproduced paths",
              "bg-error",
              drvReprod()?.UnreproducedPaths,
              drvReprod()?.Logs,
            )}
          </Show>

          <Show when={drvReprod()}>
            {renderPaths(
              "Reproduced paths",
              "bg-success",
              drvReprod()?.ReproducedPaths,
              drvReprod()?.Logs,
            )}
          </Show>

          <Show when={drvReprod()}>
            {renderPaths(
              "Unknown paths (only built by one log)",
              "bg-warning",
              drvReprod()?.UnknownPaths,
              drvReprod()?.Logs,
            )}
          </Show>

          <Show when={drvReprod()}>
            {renderPaths(
              "Missing paths (not built by any known log)",
              "bg-base-100",
              drvReprod()?.MissingPaths,
              drvReprod()?.Logs,
            )}
          </Show>
        </Suspense>
      </div>
    </>
  );
};

export default Derivation;
