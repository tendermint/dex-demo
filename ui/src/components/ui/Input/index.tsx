import React, {Component, ChangeEvent, ClipboardEvent, KeyboardEvent} from 'react';
import './input.scss';

type PropTypes = {
  className?: string
  label?: string
  placeholder?: string
  suffix?: string
  min?: string
  disabled?: boolean
  onChange?: (event: ChangeEvent<HTMLInputElement>) => void
  onKeyPress?: (event: KeyboardEvent<HTMLInputElement>) => void
  onPaste?: (event: ClipboardEvent<HTMLInputElement>) => void
  value?: string
  type?: 'text' | 'number' | 'password'
  step?: string
  autoFocus?: boolean
}

class Input extends Component<PropTypes> {
  render() {
    const {
      className = '',
      label,
      placeholder,
      disabled,
      onChange,
      onPaste,
      onKeyPress,
      suffix,
      type = 'text',
      min,
      value,
      step,
      autoFocus,
    } = this.props;

    return (
      <div className={`input ${className}`}>
        { label && <div className="input__label">{label}</div> }
        <div className="input__wrapper">
          <input
            className="input__input"
            type={type}
            placeholder={placeholder}
            disabled={disabled}
            onChange={onChange}
            onPaste={onPaste}
            onKeyPress={onKeyPress}
            min={min}
            value={value}
            step={step}
            autoFocus={autoFocus}
          />
          { suffix && <div className="input__suffix">{suffix}</div> }
        </div>
      </div>
    )
  }
}

export default Input;
