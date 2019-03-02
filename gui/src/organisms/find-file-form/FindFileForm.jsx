import React from "react";
import Input from "../../components/input/Input";
import Button from "../../components/button/Button";

import styles from "./FindFileForm.module.scss";

class FindFileForm extends React.Component {
  state = {
    filePath: ""
  };

  _onSubmit = e => {
    e.preventDefault();

    const { submit } = this.props;
    const { filePath } = this.state;

    submit(filePath);
  };

  _changeFilePath = e => {
    e.preventDefault();
    
    this.setState({ filePath: e.target.value });
  };

  render() {
    const { filePath } = this.state;
    // const { isLoading, isSuccess, error } = this.props;

    return (
      <form className={styles.FindFileForm} onSubmit={this._onSubmit}>
        <Input
          className={styles.FindFileForm__Input}
          placeholder="Enter file path..."
          value={filePath}
          onChange={this._changeFilePath}
        />
        <Button className={styles.FindFileForm__Button} type="submit">
          Find file
        </Button>
      </form>
    );
  }
}

export default FindFileForm;
