import React from "react";
import Input from "../../components/input/Input";
import Button from "../../components/button/Button";

import styles from "./FindFileForm.module.scss";

const FindFileForm = ({ onSubmit }) => (
  <form className={styles.FindFileForm} onSubmit={onSubmit}>
    <Input
      className={styles.FindFileForm__Input}
      placeholder="Enter file path..."
    />
    <Button className={styles.FindFileForm__Button} type="submit">
      Find file
    </Button>
  </form>
);

export default FindFileForm;
