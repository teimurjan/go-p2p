import React from "react";

class FindFileFormPresenter extends React.Component {
  state = {
    isLoading: false,
    error: undefined,
    isSuccess: false
  };

  _getFile = async path => {
    const response = await fetch("http://localhost:5555/getFile", {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ path })
    });

    if (response.status === 200) {
      this.setState({ isSuccess: true, isLoading: false });
    } else {
      this.setState({ error: await response.text(), isLoading: false });
    }
  };

  render() {
    const { View } = this.props;
    const { isLoading, error, isSuccess } = this.state;

    return (
      <View
        submit={this._getFile}
        isLoading={isLoading}
        error={error}
        isSuccess={isSuccess}
      />
    );
  }
}

export default FindFileFormPresenter;
