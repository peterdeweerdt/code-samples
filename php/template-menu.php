<?php
design
 */
?>

<section class="section callout-section">
  <div class="container section-content">
    <h1>
      <span class="break">
        <span class="break-mobile">From My </span>
        <span class="break-mobile">Hand </span>
      </span>
      To Your
      <span class="gold">Table.</span>
    </h1>
  </div>
</section>

<?php if($orderURL = rize_get_theme_option('rize_online_ordering_url')): ?>
  <a href="<?php echo $orderURL ?>" class="static-badge scroll-stop static-badge-order text-center"><span class="inner-text">Order<br> Online</span></a>
<?php endif; ?>

<section class="section menu-section" data-stellar-background-ratio="0.75">
  <div class="container section-content">
    <div class="content row">
      <div class="menu-gelato col-md-6">
        <div class="floating-menu-section">
          <div class="menu-img pizza"></div>
          <div class="floating-menu-wrapper first menu-1">
            <input class="accordian-trigger" id="toggle1" type="checkbox" name="toggle" />
            <label class="accordian-trigger-title" for="toggle1">
              <h2>Pizza <span class="hidden-xs">Gluten-free dough option<br/>available on all pizzas</span></h2>
            </label>

            <div class="accordian-content" id="content1">
              <div class="attention-grabber">
                <div class="attention-grabber-inner">
                  <h3>Signature <span class="menu-price">13.5</span></h3>
                  <h4>Waverly</h4>
                  <p>prosciutto, gorgonzola, shredded mozzarella, asiago, fig jam, balsamic reduction, grana padano</p>
                  <h4>Asian BBQ Chicken</h4>
                  <p>chopped bbq chicken, shredded mozzarella, charred red onion, green chili sauce, sesame seed, scallion, lime zest</p>
                  <h4>Garden</h4>
                  <p>fire roasted tomato, goat cheese, shredded mozzarella, roasted artichoke, roasted shiitake mushrooms, charred red onion, zucchini, watermelon radish, toasted pecan pesto, basil, roasted garlic oil</p>
                </div>
              </div>

              <div class="menu-5 hidden-sm hidden-md hidden-lg">
               

                <div class="row">
                  <div class="col-xs-4">
                    <h4>Cheese</h4>
                    <p>
                      asiago<br/>
                      feta<br/>
                      fresh mozzarella<br/>
                      gorgonzola<br/>
                      grana padano<br/>
                      romano<br/>
                    </p>
                  </div>
                  <div class="col-xs-4">
                    <h4>Meat</h4>
                    <p>
                      bacon<br/>
                      chopped chicken<br/>
                      crumbled fennel sausage<br/>
                      ham<br/>
                      pepperoni<br/>
                      prosciutto<br/>
                      sausage
                    </p>
                  </div>
                  <div class="col-xs-4">
                    <h4>Vegetable</h4>
                    <p>
                      caramelized onion<br/>
                      caramelized<br/>
                      pineapple<br/>
                      charred red onion<br/>
                      fire roasted tomato<br/>
                      kalamata olives<br/>
                      peppadew peppers<br/>
                      roasted artichoke<br/>
                      roasted red peppers<br/>
                      roasted shiitake<br/>
                      mushroom<br/>
                      roasted zucchini
                    </p>
                  </div>
                </div>
              </div>

              <h3>Specialty <span class="menu-price">12</span></h3>

              <h4>ATL BLT</h4>
              <p>Beeler’s bacon, baby spinach, fire roasted tomato, shredded mozzarella, charred red onion, balsamic  reduction, grana padano</p>
              <h4>Delphi</h4>
              <p>roasted artichoke, kalamata olives, fire roasted tomato, charred red onion, shredded mozzarella, feta, roasted garlic puree, baby kale</p>
              <h4>Pepperoni, Sausage & Mushroom</h4>
              <p>pepperoni, crumbled fennel sausage, roasted shiitake mushrooms, fresh mozzarella, romano, fresh oregano</p>
              <h4>Hilo</h4>
              <p>honey ham, caramelized pineapple, Beeler’s bacon, shredded mozzarella, grana padano</p>
              <h4>Peppadew & Sausage</h4>
              <p>crumbled fennel sausage, peppadew peppers, caramelized onion, fresh and shredded mozzarella, romano, fresh oregano</p>

              <h3>Classic <span class="menu-price">9.5</span></h3>

              <h4>Rosso</h4>
              <p>fire roasted tomato, fresh mozzarella, roasted garlic oil, basil, extra virgin olive oil</p>
              <h4>Bianca</h4>
              <p>roasted garlic puree, fresh and shredded mozzarella, extra virgin olive oil, basil, oregano</p>
            </div>
          </div>
        </div>

        <div class="floating-menu-section floating-menu-section-pasta hidden-sm hidden-md hidden-lg">
          <div class="menu-img pasta"></div>
          <div class="floating-menu-wrapper menu-6">
            <input class="accordian-trigger" id="toggle6" type="checkbox" name="toggle" />
            <label class="accordian-trigger-title" for="toggle6">
              <h2>Pasta <span class="menu-price hidden-xs">10</span></h2>
            </label>

            <div class="accordian-content" id="content6">
              <h4>Pork Meatball Pesto</h4>
              <p>orecchiette, pork meatballs, toasted pecan pesto, white wine, calabrese peppers, toasted bread crumb, grana padano</p>

              <h4>Roasted Shiitake Mushroom</h4>
              <p>tagliatelle pasta, roasted shiitake mushrooms, white wine, baby kale, herbed pecan bread crumbs, grana padano</p>

              <h4>Seared Shrimp</h4>
              <p>tagliatelle pasta, seared shrimp, Beeler’s bacon, fire roasted tomato, white wine & butter, parsley, grana padano</p>

              <h4>Fennel Sausage & Kale</h4>
              <p>orecchiette, crumbled fennel sausage, baby kale, fire roasted tomato, fresh mozzarella, white wine, red pepper flakes, parsley, grana padanoroasted artichoke, kalamata olives, fire roasted tomato, charred red onion, shredded mozzarella, feta, roasted garlic puree, baby kale</p>
            </div>
          </div>
        </div>

        <div class="floating-menu-section menu-2">
          <div class="menu-img salad"></div>
          <div class="floating-menu-wrapper menu-2">
            <input class="accordian-trigger" id="toggle2" type="checkbox" name="toggle" />
            <label class="accordian-trigger-title" for="toggle2">
              <h2>Salads <span class="hidden-xs">Full <span>11</span><br/>Half <span>7</span></span></h2>
            </label>

            <div class="accordian-content" id="content2">
              <input class="accordian-trigger" id="toggle2" type="checkbox" name="toggle" />
              <h4>Spice Road Chicken</h4>
              <p>sliced chicken breast, baby spinach, goat cheese, chickpeas, sweet peppers, fennel, basil, mint, cilantro, scallion, golden spice vinaigrette</p>

              <h4>Superfood</h4>
              <p>seared shrimp, quinoa, mixed lettuce, feta cheese, zucchini, cucumber, spiced walnuts, honey lemon vinaigrette</p>

              <h4>Vegetable</h4>
              <p>roasted red peppers, tri-color carrot, fire roasted tomato, zucchini, roasted artichoke, goat cheese, chickpeas, watermelon radish, fennel, baby spinach, crushed pistachio, lemon ginger vinaigrette</p>

              <h4>Pear & Walnut</h4>
              <p>baby kale & spinach, sliced pear, spiced walnuts, smoked bleu cheese, charred red onion, fennel, pomegranate seeds, honey lemon vinaigrette</p>

              <h4>Asian Shrimp</h4>
              <p>shrimp, baby kale, snap peas, cauliflower, tri-color carrot, watermelon radish, chilled soba noodles, crispy chow mein, sesame seed, sesame lime vinaigrette</p>

              <h4>Steak & Smoked Bleu Cheese Wedge  + 2</h4>
              <p>sliced angus steak, romaine, Beeler’s bacon, smoked bleu cheese, fire roasted tomato, charred red onion, bleu cheese dressing</p>
            </div>
          </div>
        </div>

        <div class="floating-menu-section hidden-sm hidden-md hidden-lg">
          <div class="menu-img hummus"></div>
          <div class="floating-menu-wrapper menu-7">
            <input class="accordian-trigger" id="toggle7" type="checkbox" name="toggle" />
            <label class="accordian-trigger-title" for="toggle7">
              <h2>Small Plates <span class="menu-price hidden-xs">7</span></h2>
            </label>

            <div class="accordian-content" id="content7">
              <h4>Old World Hummus</h4>
              <p>spiced walnuts, watermelon radish, tri color carrot, peppadew peppers, olive oil, lavash</p>

              <h4>Stuffed Peppadew Peppers</h4>
              <p>whipped feta & crumbled fennel sausage, spiced walnuts, calabrese oil</p>

              <h4>Charbroiled Wings</h4>
              <p>spice herb rub, ginger-honey sauce, cilantro</p>

              <h4>Charred Cauliflower</h4>
              <p>cauliflower florets, maple sesame dressing, spiced walnuts, pomegranate seed</p>

              <h4>Grilled Pork Meatballs</h4>
              <p>fire roasted tomato, shaved grana padano</p>

              <h4>Herbed Goat Cheese Crostini</h4>
              <p>torn flatbread crostini, pear chutney, pistachio chili crumble, balsamic reduction</p>

              <h2>Flatbreads <span class="menu-price">8</span></h2>
              <h4>Prosciutto & Pomegranate</h4>
              <p>goat cheese, baby spinach, spiced walnuts,balsamic reduction</p>

              <h4>Roasted Veggie</h4>
              <p>fire roasted tomato, artichoke, zucchini, charred red onion, fresh mozzarella, baby spinach, toasted pecan pesto, balsamic reduction</p>

              <h4>Chicken & Pear</h4>
              <p>chopped chicken breast, pear chutney, feta, charred red onion, green chili sauce</p>

              <h4>Roasted Shiitake Mushroom</h4>
              <p>caramelized onion, roasted garlic puree, baby kale, asiago, lemon zest</p>
            </div>
          </div>
        </div>

        <div class="floating-menu-section">
          <div class="menu-img sandwich"></div>
          <div class="floating-menu-wrapper menu-3">
            <input class="accordian-trigger" id="toggle3" type="checkbox" name="toggle" />
            <label class="accordian-trigger-title" for="toggle3">
              <h2>Lunch <span class="hidden-xs">11AM to 4PM <span>Daily</span></span></h2>
            </label>

            <div class="accordian-content" id="content3">
              <h3>Sandwiches <span class="menu-price">8</span></h3>
              <h4>Chicken Bacon Pesto</h4>
              <p>sliced chicken breast, Beeler’s bacon, toasted pecan pesto, fire roasted tomato, baby spinach, shredded mozzarella, asiago</p>

              <h4>Sausage & Caramelized Onion</h4>
              <p>crumbled fennel sausage, caramelized onion, roasted red peppers, fresh mozzarella, baby spinach, calabrese peppers</p>

              <h4>Meatball & Mozzarella</h4>
              <p>pork meatball, shredded mozzarella, fire roasted tomato, roasted red peppers, peppadew peppers, baby spinach</p>

              <h4>Steak & Bleu Cheese</h4>
              <p>sliced angus steak, smoked bleu cheese, caramelized onion, fire roasted tomato, charred red onion</p>

              <h3>Soup <span class="menu-price">4.5</span></h3>
              <h4>Rustic Tomato</h4>
              <h4>Seasonal Soup</h4>

              <h3>Lunch B.Y.O. Pizza <span class="menu-price">12.5</span></h3>
              <h5 class="lowercase">Includes 4 toppings and a craft beverage</h5>

              <div class="attention-grabber">
                <div class="attention-grabber-inner">
                  <h3>Pick two <span class="menu-price">10</span></h3>
                  <h5>Monday - Friday, 11 am - 2 pm only</h5>
                  <div class="pick-two">
                    <h3 class="pick-two-number">1</h3>
                    <div class="pick-two-content">
                      <h4>Choose A Flatbread Or Sandwich</h4>
                      <p>(Steak and bleu cheese sandwich)</p>
                    </div>
                    <div class="clearfix"></div>
                  </div>
                  <div class="pick-two">
                    <h3 class="pick-two-number">2</h3>
                    <div class="pick-two-content">
                      <h4>Choose A Soup Or Lunch Salad</h4>
                      <p>(Steak and bleu cheese wedge)</p>
                    </div>
                    <div class="clearfix"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="floating-menu-wrapper menu-8 hidden-sm hidden-md hidden-lg">
          <input class="accordian-trigger" id="toggle8" type="checkbox" name="toggle" />
          <label class="accordian-trigger-title" for="toggle8">
            <h2>Kids <span class="hidden-xs">Includes kid‘s beverage +<br/>kid‘s scoop (10 & under) <span>7</span></span></h2>
          </label>

          <div class="accordian-content" id="content8">
            <h4>Kid Pizza</h4>
            <p>cheese or pepperoni</p>

            <h4>Kid Pasta</h4>
            <p>orecchiette, roasted tomatoes, butter or olive oil, grana padano</p>

            <h4>Kid Hummus & Vegetables</h4>
            <p>tri colored carrot, cucumber, torn flatbread</p>
          </div>
        </div>

        <div class="floating-menu-section">
          <div class="menu-img gelato"></div>
          <div class="floating-menu-wrapper menu-4">
            <input class="accordian-trigger" id="toggle4" type="checkbox" name="toggle" />
            <label class="accordian-trigger-title" for="toggle4">
              <h2>Gelato &amp; Sorbet</h2>
            </label>

            <div class="accordian-content" id="content4">
              <h5>Single <span>3</span> Double <span>4</span> Triple <span>6</span></h5>
              <h3>Gelato</h3>
              <div class="row">
                <div class="col-md-6">
                  <p>
                    salted caramel<br/>
                    pistachio<br/>
                    tiramisu<br/>
                    cookies & cream<br/>
                    chocolate
                  </p>
                </div>
                <div class="col-md-6">
                  <p>
                    tahitian vanilla<br/>
                    bourbon brown<br/>
                    butter pecan<br/>
                    mixed berry
                  </p>
                </div>
              </div>
              <h3>Sorbet</h3>
              <p>
                strawberry<br/>
                lemon<br/>
                mango<br/>
                blood orange
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="menu-beverages col-md-6">
        <div class="floating-menu-wrapper second menu-5 hidden-xs">
          <input class="accordian-trigger" id="toggle5" type="checkbox" name="toggle" />
          <label class="accordian-trigger-title" for="toggle5">
            <h3>B.Y.O. <span class="menu-price hidden-xs">12.5</span>
              <span class="sub-title hidden-xs">Includes 4 toppings</span>
            </h3>
          </label>

          <div class="accordian-content" id="content5">
            <h5>Additional Toppings + 2.5 Each</h5>

            <div class="row">
              <div class="col-md-6">
                <h4>Cheese</h4>
                <p>
                  asiago<br/>
                  feta<br/>
                  fresh mozzarella<br/>
                  gorgonzola<br/>
                  grana padano<br/>
                  romano<br/>
                </p>

                <h4>Meat</h4>
                <p>
                  bacon<br/>
                  chopped chicken<br/>
                  crumbled fennel sausage<br/>
                  ham<br/>
                  pepperoni<br/>
                  prosciutto<br/>
                  sausage
                </p>
              </div>
              <div class="col-md-6">
                <h4>Vegetable</h4>
                <p>
                  caramelized onion<br/>
                  caramelized<br/>
                  pineapple<br/>
                  charred red onion<br/>
                  fire roasted tomato<br/>
                  kalamata olives<br/>
                  peppadew peppers<br/>
                  roasted artichoke<br/>
                  roasted red peppers<br/>
                  roasted shiitake<br/>
                  mushroom<br/>
                  roasted zucchini
                </p>
              </div>
            </div>
          </div>
        </div>

        <div class="floating-menu-section floating-menu-section-pasta hidden-xs">
          <div class="menu-img pasta"></div>
          <div class="menu-img hummus"></div>
          <div class="floating-menu-wrapper menu-6">
            <input class="accordian-trigger" id="toggle6" type="checkbox" name="toggle" />
            <label class="accordian-trigger-title" for="toggle6">
              <h2>Pasta <span class="menu-price hidden-xs">10</span></h2>
            </label>

            <div class="accordian-content" id="content6">
              <h4>Pork Meatball Pesto</h4>
              <p>orecchiette, pork meatballs, toasted pecan pesto, white wine, calabrese peppers, toasted bread crumb, grana padano</p>

              <h4>Roasted Shiitake Mushroom</h4>
              <p>bucatini, roasted shiitake mushrooms, white wine, baby kale, toasted bread crumbs, grana padano</p>

              <h4>Seared Shrimp</h4>
              <p>bucatini, seared shrimp, Beeler’s bacon, fire roasted tomato, white wine & butter, parsley, grana padano</p>

              <h4>Fennel Sausage & Kale</h4>
              <p>orecchiette, crumbled fennel sausage, baby kale, fire roasted tomato, fresh mozzarella, white wine, red pepper flakes, parsley, grana padano</p>
            </div>
          </div>
        </div>

        <div class="floating-menu-section hidden-xs">
          <div class="floating-menu-wrapper menu-7">
            <input class="accordian-trigger" id="toggle7" type="checkbox" name="toggle" />
            <label class="accordian-trigger-title" for="toggle7">
              <h2>Small Plates <span class="menu-price hidden-xs">7</span></h2>
            </label>

            <div class="accordian-content" id="content7">
              <h4>Old World Hummus</h4>
              <p>spiced walnuts, watermelon radish, tri color carrot, peppadew peppers, olive oil, lavash</p>

              <h4>Stuffed Peppadew Peppers</h4>
              <p>whipped feta & crumbled fennel sausage, spiced walnuts, calabrese oil, parsley</p>

              <h4>Charbroiled Wings</h4>
              <p>spice herb rub, ginger-honey sauce, cilantro</p>

              <h4>Charred Cauliflower</h4>
              <p>cauliflower florets, maple sesame dressing, spiced walnuts, pomegranate seed</p>

              <h4>Grilled Pork Meatballs</h4>
              <p>fire roasted tomato, shaved grana padano</p>

              <h4>Herbed Goat Cheese Crostini</h4>
              <p>torn flatbread crostini, pear chutney, pistachio chili crumble, balsamic reduction</p>

              <h2>Flatbreads <span class="menu-price">8</span></h2>
              <h4>Prosciutto & Pomegranate</h4>
              <p>goat cheese, baby spinach, spiced walnuts, balsamic reduction</p>

              <h4>Roasted Veggie</h4>
              <p>fire roasted tomato, artichoke, zucchini, charred red onion, fresh mozzarella, baby spinach, toasted pecan pesto, balsamic reduction</p>

              <h4>Chicken & Pear</h4>
              <p>sliced chicken breast, pear chutney, feta, charred red onion, green chili sauce</p>

              <h4>Roasted Shiitake Mushroom</h4>
              <p>caramelized onion, roasted garlic puree, baby kale, asiago, lemon zest</p>
            </div>
          </div>
        </div>

        <div class="floating-menu-wrapper menu-8 hidden-xs">
          <input class="accordian-trigger" id="toggle8" type="checkbox" name="toggle" />
          <label class="accordian-trigger-title" for="toggle8">
            <h2>Kids <span class="hidden-xs">Includes kid‘s beverage +<br/>kid‘s scoop (10 & under) <span>7.5</span></span></h2>
          </label>

          <div class="accordian-content" id="content8">
            <h4>Kid Pizza</h4>
            <p>cheese or pepperoni</p>

            <h4>Kid Pasta</h4>
            <p>orecchiette, roasted tomatoes, butter or olive oil, grana padano</p>

            <h4>Kid Hummus & Vegetables</h4>
            <p>tri colored carrot, cucumber, torn flatbread</p>
          </div>
        </div>

        <div class="floating-menu-section">
          <div class="menu-img sangria"></div>

          <div class="floating-menu-wrapper menu-9">
            <input class="accordian-trigger" id="toggle9" type="checkbox" name="toggle" />
            <label class="accordian-trigger-title" for="toggle9">
              <h2>Beverages</h2>
            </label>

            <div class="accordian-content" id="content9">
              <input class="accordian-trigger" id="toggle1" type="checkbox" name="toggle" />
              <label class="accordian-trigger-title" for="toggle1">
                <h3>Craft Beer <span class="menu-price">On tap 6</span></h3>
              </label>
              <p>
                Sierra Nevada (Seasonal)<br/>
                Sweetwater 420 Extra Pale Ale<br/>
                Brooklyn Lager<br/>
                Harpoon Brewery UFO White<br/>
                Lagunitas Brewing IPA<br/>
                Stiegl Radler<br/>
                Wild Heaven Emergency Drinking Beer (Local)<br/>
                Terrapin Hopsecutioner IPA
              </p>

              <h3>Wine <span class="menu-price">Glass 7</span></h3>
              <p>
                HOUSE cabernet<br/>
                NEPRICA RED IGT blend<br/>
                HOUSE chardonnay<br/>
                SAINT M riesling<br/>
                VILLA MARIA sauvignon blanc<br/>
                CIPRESSETO rose
              </p>

              <h5 class="gold">Ask about bottle and can beer selections</h5>

              <div class="attention-grabber">
                <div class="attention-grabber-inner">
                  <h3>Hand Crafted Sangria</h3>

                  <h5>Glass <span>7</span> Carafe <span>35</span></h5>

                  <h4>Marbella Red</h4>
                  <p>bold red wine infused with orange, blueberry, strawberry</p>
                  <h4>Andalucia White</h4>
                  <p>crisp white wine infused with pear, lemon, ginger</p>
                </div>
              </div>

              <h3>Craft Beverages <span class="menu-price">Bottomless</span></h3>

              <h5>Craft Soda <span>3</span></h5>
              <p>
                Caleb’s Kola<br/>
                Caleb’s Diet Kola<br />
                Agave Vanilla Cream<br/>
                Black Cherry with Tarragon<br/>
                Lemon Berry Acai<br/>
                Classic Root Beer
              </p>

              <h5>Iced Teas & Lemonade <span>2.5</span></h5>
              <p>
                Black Tea<br/>
                Sweet Black Tea<br/>
                Sweet Peach Tea<br/>
                Sweet Raspberry Tea<br/>
                Lemonade<br/>
                Strawberry Lemonade
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>
<div id="button-stop-here" class="desktop-only" style="height: 200px; margin-top: -200px;"></div>
