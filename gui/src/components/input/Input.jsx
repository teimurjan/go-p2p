import React from "react";
import classNames from "classnames";

import s from "./Input.module.scss";

const Input = ({ className, placeholder }) => (
  <div className={classNames(s.Input, className)}>
    <input type="text" className={s.Input__input} placeholder={placeholder} />
  </div>
);

export default Input;
