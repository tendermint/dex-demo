import React, { Component, ReactNode } from 'react';
import Dropdown, {ItemsType, ItemType} from '../../components/ui/Dropdown';
import './search-box.scss';

type PropTypes = {

}

type StateTypes = {
  isFocused: boolean
}

class SearchBox extends Component<PropTypes, StateTypes> {

  render () {
    return (
      <Dropdown
        items={[
          { value: '1', label: 'BTC / ETH'},
          { value: '2', label: 'BTC / USDT'},
          { value: '3', label: 'BTC / DAI'},
          { value: '4', label: 'ETH / USDT'},
          { value: '5', label: 'ETH / DAI'},
          { value: '6', label: 'ETH / BTC'},
        ]}
        currentIndex={0}
        CurrentItem={this.renderCurrentItem}
        Items={this.renderItems}
      />
    )
  }

  renderCurrentItem = (item: ItemType): ReactNode => (
    <div className="search-box">
      <div className="search-box__icon" />
      <input
        placeholder="Search by Symbol"
        onFocus={item.openDropdown}
        onBlur={item.closeDropdown}
      />
    </div>
  )

  renderItems(props: ItemsType): ReactNode {
    const { items } = props;

    const favs = items.slice(0, 2);
    const results = items.slice(2);

    return (
      <div className="dropdown__items search-box__results">
        <div className="search-box__results-group search-box__results-group--favorite">
          {
            favs.map((item: ItemType) => (
              <div className="search-box__results__item">
                <div className="search-box__results__label">
                  {item.label}
                </div>
                <div className="search-box__results__stats">
                  <div>123.47 (+2.03%)</div>
                  <div>24h Vol: 236M</div>
                </div>
              </div>
            ))
          }
        </div>
        <div className="search-box__results-group">
          {
            results.map((item: ItemType) => (
              <div className="search-box__results__item">
                <div className="search-box__results__icon" />
                <div className="search-box__results__label">
                  {item.label}
                </div>
                <div className="search-box__results__stats">
                  <div></div>
                  <div></div>
                </div>
              </div>
            ))
          }
        </div>
      </div>
    )
  }
}

export default SearchBox;
