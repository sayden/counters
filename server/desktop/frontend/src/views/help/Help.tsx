import Header from '../../components/Header';

export default function Help() {
	return (
		<main className="flex flex-col items-center">
			<div className='w-[80%] flex flex-col grow'>
				<Header className='flex justify-between items-baseline border-b-1 border-double !px-[1ch] !py-[1lh]' />
				<section>
					<h1># Counters</h1>
					<p>What is a counter?</p>
					<p>As of today, a counter is just a name that historically refer to a piece of
						a boardgame. After adding functionality, it's more like a template for that
						piece and how it is defined.</p>
					<h2> ## What does a counter look like?</h2>
					<div className='flex w-full justify-around'>
						<pre className='w-1/2'>
							<code>
								{`{
  "texts": [
    {
      "string": "Hello",
      "font_height": 25
    },
    {
      "position": 11,
      "string": "Counter",
      "font_color": "white"
    }
  ],
  "images": [
    {
      "position": 3,
      "path": "assets/supply.png",
      "scale": 0.3
    }
  ]
}
`}
							</code>
						</pre>
						<img className='object-none' src='/src/assets/images/views/help/what-is-a-counter.png' alt='what-is-a-counter' />
					</div>
					<p>
						You can see, there's 2 root elements: a <code>texts</code> and a <code>images</code>.
						Texts are a list of text definitions which requires, at least, a mandatory <code>string</code> property.
						Images are a list of image definitions and each of them requires, at least, a mandatory <code>path</code> property.
					</p>
					<img className='float-right ml-[1ch]' src='/src/assets/docs/counters/position_all.png' alt='what-is-a-counter' />
					<p>
						<b>Creating texts</b>: In the example above, the <code>texts</code> array contains 2 elements. The first text contains the text <mark>Hello</mark>,
						which you define using the <code>string</code> property. Then the size of the text is changed to <mark>25</mark> px using the optional <code>font_height</code> property.
						You can see the text <mark>HELLO</mark> in the center of the counter.
					</p>
					<p>
						The second text item, with the <code>string</code> <mark>COUNTER</mark> changes the <code>font_color</code> to white.
						But this text is in the bottom of the piece. This is define by the <code>position</code> property. And why the position is 11?
						Look at the image on the right. Each number is a position on the piece, being 0, the center; the default position (now you know
						why the text "HELLO" was positioned in the center without any <code>position</code> defined property).
					</p>
					<p>
						<b>Creating images</b>: The only image in the <code>images</code> array contains an arleady familiar <code>position</code> property on the 3rd position (at the top).
						Then, a <code>path</code> property is defined with the path in your stogare to your image. Finally, the <code>scale</code> property helps with shrinking or expanding the image.
					</p>
				</section>
				<section>
					<h2>## The settings</h2>
					<p>
						Each text and image can contain more than 25 different properties. Some of them are common, like <code>position</code> and some of them are specific like <code>font_color</code> for texts or <code>scaling</code> for images.
					</p>
					<p>
						I'll make a list with all of them, but I'll write here only most commonly used:
					</p>
					<ul>
						<li><code>position</code>: The position of the text or image in the counter. It's a number from 0 to 16. The first position is the top-left corner. The second position is the top-middle. And so on.</li>
						<li><code>font_color</code>: (Text only) The color of the text. It can be a hexadecimal color or a string color. For example, <code>font_color: "#0F0"</code> or <code>font_color: "red"</code>.</li>
						<li><code>background_color</code>: The background color of the text. It can be a hexadecimal color or a string color. For example, <code>background_color: "#0F0"</code> or <code>background_color: "red"</code>.</li>
						<li><code>x_shift</code>, <code>y_shift</code>: Displacement of the object on the counter on the x or y axis, like <code>"x_shift": 70</code></li>
						<li><code>alignment</code>: The alignment of the text. It can be <mark>center</mark>, <mark>right</mark> or <mark>left</mark>.</li>
						<li><code>image_scaling</code>: The scaling of the image. It can be <mark>fitWidth</mark>, <mark>wrap</mark> or <mark>fitHeight</mark>.</li>
					</ul>
				</section>
				<section>
					<h2>## The templates</h2>
				</section>
			</div>
		</main>
	)
}
