import {
  lazy,
  Component,
  createSignal,
  createResource,
  For,
  Show,
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


type DerivationReproducibilityPaths = { [key: string]: DerivationReproducibilityResponse_Derivation }


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
  console.log(paths)
  return (
    <>
      <p>hello</p>
      <ul>
        <For each={paths}>{(k, i) => <>
          <p>Hellolo</p>
        </>}</For>
      </ul>
    </>
  )
}

const Derivation: Component = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [drvReprod] = createResource(() => searchParams.storePath, fetchDerivationReproducibility);

  return (
    <>
      <span>{drvReprod.loading && "Loading..."}</span>

      <div>
        <h2>{searchParams.storePath}</h2>
      </div>

      <Show when={true}>
        <h3>Missing paths (not built by any known log)</h3>
        {renderPaths(drvReprod()?.MissingPaths)}
      </Show>

      <div>
        <pre>{JSON.stringify(drvReprod(), null, 2)}</pre>
      </div>
    </>
  );
};

export default Derivation;
