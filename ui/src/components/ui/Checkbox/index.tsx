import React, {ChangeEvent, Component, ReactNode} from "react";
import c from "classnames";
import "./checkbox.scss";

type Props = {
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void
  children?: ReactNode | ReactNode[]
  className?: string
  checked?: boolean
}

export default class Checkbox extends Component<Props> {
  render(): ReactNode {
    const {
      className,
      checked,
      onChange,
    } = this.props;

    return (
      <div
        className={c('checkbox', className, {
          'checkbox--clickable': !!onChange,
          'checkbox--checked': checked,
        })}
      >
        <div className="checkbox__input">
          <input
            type="checkbox"
            checked={checked}
            onChange={onChange}
            tabIndex={onChange ? 0 : -1}
            readOnly={!onChange}
          />
        </div>
        <div className="checkbox__content">
          { this.props.children }
        </div>
      </div>
    )
  }
}
