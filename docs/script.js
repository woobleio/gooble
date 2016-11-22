obj = {
  _init: function() {
    var card = this._doc.querySelector('.card');
    card.addEventListener('click', function () {
      card.classList.toggle('card--open');
    });
  }
}
