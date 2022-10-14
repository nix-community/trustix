import {
  Component,
  createResource,
  For,
  Show,
  Suspense,
  createEffect,
} from "solid-js";
import { useSearchParams, A } from "@solidjs/router";
import { createStore } from "solid-js/store";

import {
  DerivationReproducibilityRequest,
  DerivationReproducibilityResponse,
  DerivationReproducibilityResponse_Derivation,
  DerivationReproducibilityResponse_DerivationOutput,
  DerivationReproducibilityResponse_DerivationOutputHash,
} from "../api/api_pb";
import { NameValuePair } from "../lib";
import client from "../client";

import { loading } from "../widgets";

import { SolidChart, SolidChartProps } from "../chart/SolidChart";

type DerivationReproducibilityPaths = {
  [key: string]: DerivationReproducibilityResponse_Derivation;
};
type Logs = { [key: string]: Log };

/* eslint-disable @typescript-eslint/no-explicit-any */
const objSize = (o: any): number => Object.keys(o).length;

const fetchDerivationReproducibility = async (
  drvPath,
): DerivationReproducibilityResponse => {
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
                class="checkbox"
                {...checkboxArgs}
              />
            </label>
          </th>
          <td>
            <div class="flex items-center space-x-3">
              <div>
                <div class="text-sm opacity-50">{outputNarinfoHash}</div>
              </div>
            </div>
          </td>
          <td>
            <For each={logIDs}>
              {(logID) => (
                <>
                  <span class="badge badge-ghost badge-sm">
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

  const onNarinfoClicked = () => {
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

    alert("TODO: Redirect to diff view: ", a, b);
  };

  return (
    <>
      <div class="card bg-base-200 shadow-xl m-3">
        <div class="card-body">
          <h2 class="card-title tooltip" data-tip="Output name">
            {output}
          </h2>
          <p class="font-bold">{storePath}</p>

          {outputHashes.length > 0 && (
            <>
              <div class="overflow-x-auto w-full">
                <table class="table w-full">
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
                    class="btn btn-info btn-sm"
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
      <div class={`card drv-card shadow-xl m-2 ${cardBackground}`}>
        <div class="card-body">
          <A
            href={`/drv?storePath=${encodeURIComponent(
              encodeURIComponent(drvPath),
            )}`}
          >
            <h2 class="card-title tooltip" data-tip="Derivation store path">
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
  const derivations = NameValuePair.fromMap(paths);

  return (
    <>
      <div class="divider" />

      <div class="grid place-items-center w-11/12">
        <h3 class="text-xl font-bold underline">{title}</h3>

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

const renderDerivationStatistics = (
  drvReprod: DerivationReproducibilityResponse,
): Component => {
  // Show a doughnut chart with the different groups
  const chartSettings: SolidChartProps = {
    type: "doughnut",
    data: {
      labels: [
        "Unreproduced paths",
        "Reproduced paths",
        "Unknown paths",
        "Missing paths",
      ],
    },
    options: {
      responsive: false,
      plugins: {
        legend: {
          position: "top",
        },
      },
    },
  };

  const [chart, setChart] = createStore(chartSettings);

  createEffect(() => {
    const resp = drvReprod;
    if (resp == undefined) {
      return;
    }

    const datasets = [
      {
        data: [
          resp.UnreproducedPaths,
          resp.ReproducedPaths,
          resp.UnknownPaths,
          resp.MissingPaths,
        ].map((paths) => Object.keys(paths).length),
        backgroundColor: [
          "#f87272", // bg-error
          "#36d399", // bg-success
          "#fbbd23", // bg-warning
          "#ffffff", // bg-base-100
        ],
      },
    ];

    setChart("data", "datasets", datasets);
  });

  /* eslint-disable solid/reactivity */
  const numOutputs = [
    drvReprod.UnreproducedPaths,
    drvReprod.ReproducedPaths,
    drvReprod.UnknownPaths,
    drvReprod.MissingPaths,
  ].reduce((acc, v) => acc + objSize(v), 0);

  const numReproduced = objSize(drvReprod.ReproducedPaths);

  return (
    <div class="flex justify-evenly">
      <div class="card w-96 bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">Statistics</h2>

          <table class="table">
            <tbody>
              <tr>
                <td class="font-bold">Unreproduced paths</td>
                <td>{objSize(drvReprod.UnreproducedPaths)}</td>
              </tr>

              <tr>
                <td class="font-bold">Reproduced paths</td>
                <td>{numReproduced}</td>
              </tr>

              <tr>
                <td class="font-bold">Unknown paths</td>
                <td>{objSize(drvReprod.UnknownPaths)}</td>
              </tr>

              <tr>
                <td class="font-bold">Missing paths</td>
                <td>{objSize(drvReprod.MissingPaths)}</td>
              </tr>

              <tr>
                <td class="font-bold">Reproduced</td>
                <td>{(numOutputs / 100) * numReproduced}%</td>
              </tr>

              <tr>
                <td class="font-bold">Number of logs</td>
                <td>{objSize(drvReprod.Logs)}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div>
        <SolidChart
          {...chart}
          canvasOptions={{
            width: 300,
            height: 300,
          }}
        />
      </div>
    </div>
  );
};

const Derivation: Component = () => {
  const [searchParams] = useSearchParams();
  const [drvReprod] = createResource(
    () => searchParams.storePath,
    fetchDerivationReproducibility,
  );

  return (
    <>
      <div>
        <h2 class="text-xl font-bold text-center mb-2">
          {searchParams.storePath}
        </h2>

        <Suspense fallback={loading}>
          <Show when={drvReprod()}>
            {renderDerivationStatistics(drvReprod())}
          </Show>

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
