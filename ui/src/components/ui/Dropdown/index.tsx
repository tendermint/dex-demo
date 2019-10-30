import React, { Component, ReactNode } from 'react';
import c from 'classnames';
import './dropdown.scss';

export type ItemType = {
  label: string
  value: string
  toggleDropdown?: () => void
  openDropdown?: () => void
  closeDropdown?: () => void
  meta?: any
}

export type ItemsType = {
  items: ItemType[]
  toggleDropdown?: () => void
  openDropdown?: () => void
  closeDropdown?: () => void
}

type PropTypes = {
  className?: string
  items: ItemType[]
  currentIndex: number
  Items?: (props: ItemsType) => ReactNode | (new (props: ItemsType) => Component)
  Item?: (props: ItemType) => ReactNode | (new (props: ItemType) => Component)
  CurrentItem?: (props: ItemType) => ReactNode | (new (props: ItemType) => Component)
  onSelect?: (val: string, index: number) => void
}

type StateTypes = {
  isOpened: boolean
}

export default class Dropdown extends Component<PropTypes, StateTypes> {
  state = {
    isOpened: false,
  };

  private extendItemProps(props: ItemType): ItemType {
    return {
      ...props,
      openDropdown: this.openDropdown,
      closeDropdown: this.closeDropdown,
      toggleDropdown: this.toggleDropdown,
    }
  }

  private extendItemsProps(props: ItemsType): ItemsType {
    return {
      ...props,
      openDropdown: this.openDropdown,
      closeDropdown: this.closeDropdown,
      toggleDropdown: this.toggleDropdown,
    }
  }

  openDropdown = () => this.setState({ isOpened: true });
  closeDropdown = () => this.setState({ isOpened: false });

  toggleDropdown = () => this.setState({
    isOpened: !this.state.isOpened,
  });

  selectItem = (val: string, index: number) => {
    const { onSelect } = this.props;
    if (typeof onSelect === 'function') onSelect(val, index);
    this.toggleDropdown();
  }

  render() {
    const {
      className = '',
    } = this.props;

    const { isOpened } = this.state;

    return (
      <div
        className={c(`dropdown ${className}`, {
          'dropdown--opened': isOpened,
          'dropdown--closed': !isOpened,
        })}
      >
        { isOpened && this.renderOverlay() }
        { this.renderCurrentItem() }
        { isOpened && this.renderItems() }
      </div>
    );
  }

  renderCurrentItem= () => {
    const { currentIndex, items, CurrentItem } = this.props;
    let currentItem = items[currentIndex];

    if (!currentItem) {
      currentItem = { label: ' - ', value: '' };
    }

    return CurrentItem
      ? CurrentItem(this.extendItemProps(currentItem))
      : (
        <button
          className="dropdown__item dropdown__item--current"
          onClick={this.toggleDropdown}
        >
          {currentItem.label}
        </button>
      )
  };

  renderItems = (): React.ReactNode => {
    const {
      items,
      Items,
    } = this.props;

    return Items
      ? Items(this.extendItemsProps({ items }))
      : (
        <div className="dropdown__items">
          {items.map(this.renderItem)}
        </div>
      )
  }

  renderItem = (props: ItemType, i: number): React.ReactNode => {
    const { Item } = this.props;

    return Item
      ? Item(this.extendItemProps(props))
      : (
        <div
          key={props.value}
          className="dropdown__item"
          onClick={() => this.selectItem(props.value, i)}
        >
          {props.label}
        </div>
      )
  }

  renderOverlay(): React.ReactNode {
    return (
      <div className="dropdown__overlay" onClick={this.toggleDropdown}/>
    );
  }
}
