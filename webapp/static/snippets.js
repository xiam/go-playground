'use strict';

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

(function() {

  function bindToggle(el) {
    $('.toggleButton', el).click(function() {
      if ($(el).is('.toggle')) {
        $(el).addClass('toggleVisible').removeClass('toggle');
      } else {
        $(el).addClass('toggle').removeClass('toggleVisible');
      }
    });
    if ($(el).find('.code').data("expanded")) {
      $(el).addClass('toggleVisible').removeClass('toggle');
    };
  }

  function bindToggles(selector) {
    $(selector).each(function(i, el) {
      bindToggle(el);
    });
  }

  function setupInlinePlayground() {
    'use strict';
    // Set up playground when each element is toggled.
    $('div.play').each(function (i, el) {
      // Set up playground for this example.
      var setup = function() {
        var code = $('.code', el);

        playground({
          'codeEl':   code,
          'outputEl': $('.output', el),
          'runEl':    $('.run', el),
          'fmtEl':    $('.fmt', el),
          'shareEl':  $('.share', el),
        });
      };

      // If example already visible, set up playground now.
      if ($(el).is(':visible')) {
        setup();
        return;
      }

      // Otherwise, set up playground when example is expanded.
      var built = false;
      $(el).closest('.toggle').click(function() {
        // Only set up once.
        if (!built) {
          setup();
          built = true;
        }
      });
    });
  }

  function initGoPlayground() {
    var snippets = $('textarea.go-playground-snippet');
    for (var i = 0; i < snippets.length; i++) {
      var el = $(snippets[i]);
      var title = el.data("title") || "Snippet";
      var init = el.data("initialized");
      if (!init) {
        var html = ''+
        '<div class="toggle">'+
        ' <div class="collapsed">'+
        '   <p class="exampleHeading toggleButton">▹ <span class="text">Example</span></p>'+
        ' </div>'+
        ' <div class="expanded">'+
        '   <p class="exampleHeading toggleButton">▾ <span class="text">Example</span></p>'+
        '     <div class="play">'+
        '       <div class="input"></div>'+
        '       <div class="output"></div>' +
        '       <div class="buttons">' +
        '         <a class="run" title="Run this code [shift-enter]">Run</a>' +
        '         <a class="fmt" title="Format this code">Format</a>' +
        '         <a class="share" title="Share this code">Share</a>' +
        '       </div>' +
        '     </div>' +
        ' </div>' +
        '</div>'
        var wrapper = $(html);
        wrapper.find('p.exampleHeading .text').text(title);
        wrapper.insertBefore(el);
        wrapper.find('div.input').append(el);
        el.addClass("code");
        el.attr("autocorrect", "off");
        el.attr("autocomplete", "off");
        el.attr("autocapitalize", "off");
        el.attr("spellcheck", "false");
        el.data("initialized", 1);
      }
    }
  }

  $(document).ready(function() {
    initGoPlayground();
    bindToggles(".toggle");
    setupInlinePlayground();
  });
})();
