import React, { Component, ReactNode } from 'react';
import c from 'classnames';
import Dropdown, { ItemType } from '../../components/ui/Dropdown';
import './network-dropdown.scss';

type PropTypes = {
  items: ItemType[]
  // Item?: new (props: ItemType) => React.Component;
}

class NetworkDropdown extends Component<PropTypes> {
  render() {
    const { items } = this.props;

    return (
      <Dropdown
        items={items}
        currentIndex={0}
        className="network-dropdown"
        CurrentItem={(item: ItemType): ReactNode => (
          <div
            className="network-dropdown__selected-item"
            onClick={item.toggleDropdown}
          >
            <div className="network-dropdown__caret" />
            <div className="network-dropdown__selected-item__label">
              { item.label }
            </div>
            <div className="network-dropdown__selected-item__content">
              { this.renderSyncStatus(item.meta.sync) }
            </div>
          </div>
        )}
        Item={(item: ItemType): ReactNode => {
          return (
            <div
              className="network-dropdown__selected-item network-dropdown__item"
              onClick={item.toggleDropdown}
            >
              <div className="network-dropdown__selected-item__label">
                { item.label }
              </div>
              <div className="network-dropdown__selected-item__content">
                { this.renderSyncStatus(item.meta.sync) }
              </div>
            </div>
          );
        }}
      />
    );
  }

  renderSyncStatus (sync: number): ReactNode {
    return (
      <div
        className={c("network-dropdown__sync-status", {
          'network-dropdown__sync-status--synced': sync === 1,
          'network-dropdown__sync-status--syncing': sync < 1,
          'network-dropdown__sync-status--error': false,
        })}
      >
        <div
          className="network-dropdown__sync-status__percentage"
          style={{
            width: `${sync * 100}%`,
          }}
        />
        {
          sync === 1
            ? 'Synced'
            : `${(sync * 100).toFixed(0)}%`
        }
      </div>
    );
  }

  // renderLatency (latency: number): ReactNode {
  //   const text = latency < 1000
  //     ? `${latency}ms`
  //     : `${(latency/1000).toFixed(2)}s`;
  //
  //   return (
  //     <div
  //       className={c('network-dropdown__latency', {
  //         'network-dropdown__latency--fast': latency <= 250,
  //         'network-dropdown__latency--normal': latency > 250 && latency <= 500,
  //         'network-dropdown__latency--slow': latency > 500,
  //       })}
  //     >
  //       {text}
  //     </div>
  //   );
  // }
}

export default NetworkDropdown;
