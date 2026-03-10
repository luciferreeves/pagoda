import { Suspense } from "solid-js";
import type { JSX } from "solid-js";
import Layout from "./components/Layout";

const App = (props: { children?: JSX.Element }) => {
  return (
    <Layout>
      <Suspense>{props.children}</Suspense>
    </Layout>
  );
};

export default App;