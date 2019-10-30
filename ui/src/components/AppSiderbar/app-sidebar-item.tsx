import React, { Component, MouseEvent } from 'react';
import c from 'classnames';

type Props = {
  imageUrl?: String
  selected?: Boolean
  hoverable?: Boolean
  onClick?: (event: MouseEvent) => void
}

class AppSidebarItem extends Component<Props> {
  render() {
    const {
      imageUrl,
      selected,
      hoverable,
      onClick,
    } = this.props;

    return (
      <div
        className={c('app-sidebar__item', {
          'app-sidebar__item--active': selected,
          'app-sidebar__item--hoverable': hoverable,
        })}
        style={{
          backgroundImage: imageUrl ? `url(${imageUrl})` : '',
        }}
        onClick={onClick}
      />
    )
  }
}

export default AppSidebarItem;
