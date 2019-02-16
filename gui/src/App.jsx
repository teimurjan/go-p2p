import React, { Component } from "react";
import Container from "./components/container/Container";
import FindFileForm from "./organisms/find-file-form/FindFileForm";
import Content from "./components/content/Content";

class App extends Component {
  render() {
    return (
      <Container>
        <Content>
          <FindFileForm />
        </Content>
      </Container>
    );
  }
}

export default App;
