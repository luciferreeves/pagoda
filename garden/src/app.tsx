import { Suspense, type Component } from "solid-js";
import Layout from "./components/Layout";

const App: Component<{ children: Element }> = (props) => {
  return (
    <Layout>
      <Suspense>{props.children}</Suspense>
    </Layout>
  );
};

export default App;