import {
  lazy,
  Component,
  createSignal,
  createResource,
  For,
  Show,
  Suspense,
} from "solid-js";
import {
  Routes,
  Route,
  useParams,
  useSearchParams,
} from "@solidjs/router";

import {
  createConnectTransport,
  createPromiseClient,
} from '@bufbuild/connect-web'


import { ReproducibilityAPI } from '../api/api_connectweb'
import {
  DerivationReproducibilityRequest,
  DerivationReproducibilityResponse,
  DerivationReproducibilityResponse_Derivation,
} from '../api/api_pb'
import {
  NameValuePair,
} from '../lib'


type DerivationReproducibilityPaths = { [key: string]: DerivationReproducibilityResponse_Derivation }

const loading = (
  <h1>Loading...</h1>
)

const fetchDerivationReproducibility = async (drvPath): DerivationReproducibilityResponse => {
  const client = createPromiseClient(
    ReproducibilityAPI,
    createConnectTransport({
      baseUrl: '/api',
    })
  )

  const req = new DerivationReproducibilityRequest({
    DrvPath: drvPath,
  })

  return await client.derivationReproducibility(req)
}

const renderPaths = (paths: DerivationReproducibilityPaths): Component => {
  const derivations = NameValuePair.fromMap(paths)

  return (
    <ul>
      <For each={derivations}>{({name, value}) => {
          const drvPath = name
          const drvOutputs = NameValuePair.fromMap(value.Outputs)

          return (
            <div>
              <h4>{drvPath}</h4>
              <ul>
                <For each={drvOutputs}>{({name, value}) => {
                    const output = name
                    const storePath = value.StorePath

                    const outputHashes = NameValuePair.fromMap(value.OutputHashes)

                    console.log(outputHashes)

                    return (
                      <li>
                        <p>{output}</p>
                      </li>
                    )
                  }}</For>
              </ul>
            </div>
          )
        }}</For>
    </ul>
  )
}

const Derivation: Component = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [drvReprod] = createResource(() => searchParams.storePath, fetchDerivationReproducibility);

  return (
    <>
      <Suspense fallback={loading}>
        <div>
          <h2>{searchParams.storePath}</h2>
        </div>

        <Show when={drvReprod()?.UnreproducedPaths}>
          <h3>Unreproduced paths</h3>
          {renderPaths(drvReprod()?.UnreproducedPaths)}
        </Show>

        <Show when={drvReprod()?.ReproducedPaths}>
          <h3>Reproduced paths</h3>
          {renderPaths(drvReprod()?.ReproducedPaths)}
        </Show>

        <Show when={drvReprod()?.UnknownPaths}>
          <h3>Unknown paths (only built by one log)</h3>
          {renderPaths(drvReprod()?.UnknownPaths)}
        </Show>

        <Show when={drvReprod()?.MissingPaths}>
          <h3>Missing paths (not built by any known log)</h3>
          {renderPaths(drvReprod()?.MissingPaths)}
        </Show>

        <div>
          <pre>{JSON.stringify(drvReprod(), null, 2)}</pre>
        </div>
      </Suspense>
    </>
  );
};

export default Derivation;
