import React from "react";
import classNames from "classnames";

import styles from "./Button.module.scss";

const Button = ({ children, className, ...props }) => (
  <button className={classNames(styles.Button, className)} {...props}>
    {children}
  </button>
);

export default Button;
