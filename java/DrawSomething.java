import javafx.application.*;
import javafx.event.ActionEvent;
import javafx.geometry.*;
import javafx.scene.*;
import javafx.scene.control.*;
import javafx.scene.input.*;
import javafx.scene.layout.*;
import javafx.scene.paint.*;
import javafx.scene.text.*;
import javafx.scene.shape.*;
import javafx.stage.*;

enum ShapeType {
	CIRCLE_SHAPE, RECTANGLE_SHAPE, ELLIPSE_SHAPE;
}

enum PenStatus {
	DRAW, OFF, ERASE;
}

public class DrawSomething extends Application {

	private final static int DEFAULT_PADDING = 5, DEFAULT_SPACING = 30;
	private final static int SMALL_FONT_SIZE = 14;
	private final static int SCENE_WIDTH = 1000, SCENE_HEIGHT = 500;
	private final static int SMALL_BRUSH = 5, MEDIUM_BRUSH = 20, LARGE_BRUSH = 30;

	// declare instance objects
	private boolean isPenDown;
	private BorderPane borderPane;
	private Button clearButton, eraseButton;
	private CheckBox fillCheckBox;
	private Circle circle;
	private Color currentColor, previousColor;
	private Ellipse ellipse;
	private HBox pickColorHBox, penStatusHBox;
	private int brushSize;
	private Pane drawPane;
	private PenStatus penStatus;
	private RadioButton smallBrushRadioButton, mediumBrushRadioButton, largeBrushRadioButton;
	private RadioButton redRadioButton, greenRadioButton, blueRadioButton;
	private RadioButton circleBrushRadioButton, rectangleBrushRadioButton, ellipseBrushRadioButton;
	private Rectangle rectangle;
	private ShapeType currentShape;
	private Text pickBrushShapeVBoxText, pickBrushSizeVBoxText, pickBrushColorHBoxText;
	private Text penStatusHBoxText, pickShapeFillText;
	private VBox pickBrushShapeVBox, pickBrushSizeVBox, pickShapeFillVBox;

	public void start(Stage primaryStage) {

		// initialize and format borderPane
		borderPane = new BorderPane();
		borderPane.setPadding(new Insets(DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING));
		borderPane.setStyle("-fx-background-color: white");

		setupDrawPane();
		setupBrushColor();
		setupPenStatusHBox();
		setupBrushShape();
		setupShapeFill();
		setupBrushSize();

		// initialize shape objects
		circle = new Circle();
		rectangle = new Rectangle();
		ellipse = new Ellipse();

		penStatus = PenStatus.DRAW;
		isPenDown = true;

		// penStatusVBox button handlers
		clearButton.setOnAction(this::handleClearButton);

		// pickBrushShapeVBox buttons handlers
		circleBrushRadioButton.setOnAction(this::handleBrushShapeButtons);
		rectangleBrushRadioButton.setOnAction(this::handleBrushShapeButtons);
		ellipseBrushRadioButton.setOnAction(this::handleBrushShapeButtons);

		// pickBrushSizeVBox buttons handlers
		smallBrushRadioButton.setOnAction(this::handleBrushSizeButtons);
		mediumBrushRadioButton.setOnAction(this::handleBrushSizeButtons);
		largeBrushRadioButton.setOnAction(this::handleBrushSizeButtons);

		// pickBrushColorHBox buttons handlers
		redRadioButton.setOnAction(this::handleBrushColorButtons);
		greenRadioButton.setOnAction(this::handleBrushColorButtons);
		blueRadioButton.setOnAction(this::handleBrushColorButtons);

		// eraseButton handler
		eraseButton.setOnAction(this::handleEraseButton);

		// drawHBox mouse handlers
		drawPane.setOnMouseMoved(this::handleMouseMotion);
		drawPane.setOnMouseClicked(this::handleMouseClicks);

		// setup and display scene
		Scene scene = new Scene(borderPane, SCENE_WIDTH, SCENE_HEIGHT, Color.WHITE);
		primaryStage.setTitle("draw something");
		primaryStage.setScene(scene);
		primaryStage.show();

	}

	private void setupPenStatusHBox() {

		// initialize and format penStatusVBox
		penStatusHBox = new HBox();
		penStatusHBox.setAlignment(Pos.CENTER);
		penStatusHBox.setPadding(new Insets(DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING));
		penStatusHBox.setSpacing(DEFAULT_SPACING);

		// initialize and format penStatusVBoxText Object
		penStatusHBoxText = new Text();
		penStatusHBoxText.setText("draw");
		penStatusHBoxText.setFont(Font.font("Helvetica", SMALL_FONT_SIZE));
		penStatusHBoxText.setTextAlignment(TextAlignment.CENTER);

		// add penStatusHBoxText to penStatusHBox
		penStatusHBox.getChildren().add(penStatusHBoxText);

		// add clearButton to penStatusVBox
		clearButton = new Button("clear");
		clearButton.setFont(Font.font(SMALL_FONT_SIZE));
		penStatusHBox.getChildren().add(clearButton);

		// add penStatusVBox to the top borderPane
		borderPane.setTop(penStatusHBox);
	}

	private void setupBrushShape() {

		// initialize and format pickBrushShapeControlVBox
		pickBrushShapeVBox = new VBox();
		pickBrushShapeVBox.setAlignment(Pos.TOP_LEFT);
		pickBrushShapeVBox.setPadding(new Insets(DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING));
		pickBrushShapeVBox.setSpacing(DEFAULT_SPACING);

		// initialize and format pickBrushSizeControlVBoxText Object
		pickBrushShapeVBoxText = new Text();
		pickBrushShapeVBoxText.setText("brush shape");
		pickBrushShapeVBoxText.setFont(Font.font("Helvetica", SMALL_FONT_SIZE));
		pickBrushShapeVBoxText.setTextAlignment(TextAlignment.CENTER);

		// add pickBrushSizeControlVBoxText to pickBrushSizeControlVBox
		pickBrushShapeVBox.getChildren().add(pickBrushShapeVBoxText);

		// setup selectBrushSizeButtonGroup ToggleGroup
		ToggleGroup selectBrushShapeButtonGroup = new ToggleGroup();

		// setup the smallBrushRadioButton
		circleBrushRadioButton = new RadioButton("circle");
		circleBrushRadioButton.setFont(Font.font(SMALL_FONT_SIZE));
		circleBrushRadioButton.setSelected(true);

		currentShape = ShapeType.CIRCLE_SHAPE;

		circleBrushRadioButton.setToggleGroup(selectBrushShapeButtonGroup);
		pickBrushShapeVBox.getChildren().add(circleBrushRadioButton);

		// setup the mediumBrushRadioButton
		rectangleBrushRadioButton = new RadioButton("rectangle");
		rectangleBrushRadioButton.setFont(Font.font(SMALL_FONT_SIZE));
		rectangleBrushRadioButton.setToggleGroup(selectBrushShapeButtonGroup);
		pickBrushShapeVBox.getChildren().add(rectangleBrushRadioButton);

		// setup the largeBrushRadioButton
		ellipseBrushRadioButton = new RadioButton("ellipse");
		ellipseBrushRadioButton.setFont(Font.font(SMALL_FONT_SIZE));
		ellipseBrushRadioButton.setToggleGroup(selectBrushShapeButtonGroup);
		pickBrushShapeVBox.getChildren().add(ellipseBrushRadioButton);

		// add pickBrushShapeControlVBox to the left borderPane
		borderPane.setLeft(pickBrushShapeVBox);

	}

	private void setupBrushSize() {

		// initialize and format pickBrushSizeControlVBox
		pickBrushSizeVBox = new VBox();
		pickBrushSizeVBox.setAlignment(Pos.TOP_LEFT);
		pickBrushSizeVBox.setPadding(new Insets(DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING));
		pickBrushSizeVBox.setSpacing(DEFAULT_SPACING);

		// initialize and format pickBrushSizeControlVBoxText Object
		pickBrushSizeVBoxText = new Text();
		pickBrushSizeVBoxText.setText("brush size");
		pickBrushSizeVBoxText.setFont(Font.font("Helvetica", SMALL_FONT_SIZE));
		pickBrushSizeVBoxText.setTextAlignment(TextAlignment.CENTER);

		// add pickBrushSizeControlVBoxText to pickBrushSizeControlVBox
		pickBrushSizeVBox.getChildren().add(pickBrushSizeVBoxText);

		// setup selectBrushSizeButtonGroup ToggleGroup
		ToggleGroup selectBrushSizeButtonGroup = new ToggleGroup();

		// setup the smallBrushRadioButton
		smallBrushRadioButton = new RadioButton("small");
		smallBrushRadioButton.setFont(Font.font(SMALL_FONT_SIZE));
		smallBrushRadioButton.setSelected(true);
		brushSize = SMALL_BRUSH;
		smallBrushRadioButton.setToggleGroup(selectBrushSizeButtonGroup);
		pickBrushSizeVBox.getChildren().add(smallBrushRadioButton);

		// setup the mediumBrushRadioButton
		mediumBrushRadioButton = new RadioButton("medium");
		mediumBrushRadioButton.setFont(Font.font(SMALL_FONT_SIZE));
		mediumBrushRadioButton.setToggleGroup(selectBrushSizeButtonGroup);
		pickBrushSizeVBox.getChildren().add(mediumBrushRadioButton);

		// setup the largeBrushRadioButton
		largeBrushRadioButton = new RadioButton("large");
		largeBrushRadioButton.setFont(Font.font(SMALL_FONT_SIZE));
		largeBrushRadioButton.setToggleGroup(selectBrushSizeButtonGroup);
		pickBrushSizeVBox.getChildren().add(largeBrushRadioButton);

		// add pickBrushSizeControlVBox to the right borderPane
		borderPane.setRight(pickBrushSizeVBox);
	}

	private void setupShapeFill() {

		// initialize and format pickShapeFillVBox
		pickShapeFillVBox = new VBox();
		pickShapeFillVBox.setAlignment(Pos.TOP_LEFT);
		pickShapeFillVBox.setPadding(new Insets(DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING));
		pickShapeFillVBox.setSpacing(DEFAULT_SPACING);

		// initialize and format pickShapeFillText Object
		pickShapeFillText = new Text();
		pickShapeFillText.setText("brush fill");
		pickShapeFillText.setFont(Font.font("Helvetica", SMALL_FONT_SIZE));
		pickShapeFillText.setTextAlignment(TextAlignment.CENTER);

		// add pickShapeFillText to pickShapeFillVBox
		pickShapeFillVBox.getChildren().add(pickShapeFillText);

		// setup the fillCheckBox
		fillCheckBox = new CheckBox("filled?");

		fillCheckBox.setFont(Font.font(SMALL_FONT_SIZE));
		fillCheckBox.setSelected(false);
		pickShapeFillVBox.getChildren().add(fillCheckBox);
		pickBrushShapeVBox.getChildren().add(pickShapeFillVBox);

	}

	private void setupBrushColor() {

		// initialize and format pickColorControlHBox
		pickColorHBox = new HBox();
		pickColorHBox.setAlignment(Pos.CENTER);
		pickColorHBox.setPadding(new Insets(DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING, DEFAULT_PADDING));
		pickColorHBox.setSpacing(DEFAULT_SPACING);

		// initialize and format pickColorControlHBoxText Object
		pickBrushColorHBoxText = new Text();
		pickBrushColorHBoxText.setText("brush color");
		pickBrushColorHBoxText.setFont(Font.font("Helvetica", SMALL_FONT_SIZE));
		pickBrushColorHBoxText.setTextAlignment(TextAlignment.CENTER);

		// add pickColorControlHBoxText to pickColorControlHBox
		pickColorHBox.getChildren().add(pickBrushColorHBoxText);

		// setup selectColorButtonGroup ToggleGroup
		ToggleGroup selectBrushColorButtonGroup = new ToggleGroup();

		// setup the redRadioButton
		redRadioButton = new RadioButton("red");
		redRadioButton.setFont(Font.font(SMALL_FONT_SIZE));
		redRadioButton.setSelected(true);
		currentColor = Color.RED;
		previousColor = Color.RED;
		redRadioButton.setToggleGroup(selectBrushColorButtonGroup);
		pickColorHBox.getChildren().add(redRadioButton);

		// setup the yellowRadioButton
		greenRadioButton = new RadioButton("green");
		greenRadioButton.setFont(Font.font(SMALL_FONT_SIZE));
		greenRadioButton.setToggleGroup(selectBrushColorButtonGroup);
		pickColorHBox.getChildren().add(greenRadioButton);

		// setup the blueRadioButton
		blueRadioButton = new RadioButton("blue");
		blueRadioButton.setFont(Font.font(SMALL_FONT_SIZE));
		blueRadioButton.setToggleGroup(selectBrushColorButtonGroup);
		pickColorHBox.getChildren().add(blueRadioButton);

		// setup the eraseRadioButton
		eraseButton = new Button("eraser");
		eraseButton.setFont(Font.font(SMALL_FONT_SIZE));
		pickColorHBox.getChildren().add(eraseButton);

		// add pickColorControlHBox to the bottom borderPane
		borderPane.setBottom(pickColorHBox);
	}

	private void setupDrawPane() {

		// initialize and format drawHBox
		drawPane = new Pane();

		drawPane.setStyle("-fx-background-color: cyan");

		// add drawPane to the center borderPane
		borderPane.setCenter(drawPane);

	}

	// button and mouse handlers

	private void handleClearButton(ActionEvent event) {

		drawPane.getChildren().clear();

	}

	// pickBrushShapeVBox buttons handlers
	private void handleBrushShapeButtons(ActionEvent event) {

		if (circleBrushRadioButton.isSelected()) {
			currentShape = ShapeType.CIRCLE_SHAPE;
		} else if (rectangleBrushRadioButton.isSelected()) {
			currentShape = ShapeType.RECTANGLE_SHAPE;
		} else if (ellipseBrushRadioButton.isSelected()) {
			currentShape = ShapeType.ELLIPSE_SHAPE;
		}
	}

	// pickBrushSizeVBox buttons handlers
	private void handleBrushSizeButtons(ActionEvent event) {

		if (smallBrushRadioButton.isSelected()) {
			brushSize = SMALL_BRUSH;
		} else if (mediumBrushRadioButton.isSelected()) {
			brushSize = MEDIUM_BRUSH;
		} else if (largeBrushRadioButton.isSelected()) {
			brushSize = LARGE_BRUSH;
		}
	}

	// pickBrushColorHBox buttons handlers
	private void handleBrushColorButtons(ActionEvent event) {

		if (redRadioButton.isSelected()) {
			previousColor = Color.RED;
			currentColor = Color.RED;
		} else if (greenRadioButton.isSelected()) {
			previousColor = Color.GREEN;
			currentColor = Color.GREEN;
		} else if (blueRadioButton.isSelected()) {
			previousColor = Color.BLUE;
			currentColor = Color.BLUE;
		}
	}

	private void handleEraseButton(ActionEvent event) {

		penStatus = PenStatus.ERASE;
		penStatusHBoxText.setText("erase");
		currentColor = Color.CYAN;

	}

	private void handleMouseMotion(MouseEvent event) {

		double x = event.getX();
		double y = event.getY();

		if (penStatus == PenStatus.DRAW) {

			if (currentShape == ShapeType.CIRCLE_SHAPE) {
				circle = new Circle(x, y, brushSize);
				circle.setStroke(currentColor);

				if (fillCheckBox.isSelected()) {
					circle.setFill(currentColor);
				} else {
					circle.setFill(null);
				}

				drawPane.getChildren().add(circle);

			} else if (currentShape == ShapeType.RECTANGLE_SHAPE) {
				rectangle = new Rectangle(x, y, 1.3 * brushSize, brushSize);
				rectangle.setStroke(currentColor);

				if (fillCheckBox.isSelected()) {
					rectangle.setFill(currentColor);
				} else {
					rectangle.setFill(null);
				}

				drawPane.getChildren().add(rectangle);

			} else if (currentShape == ShapeType.ELLIPSE_SHAPE) {
				ellipse = new Ellipse(x, y, 1.3 * brushSize, brushSize);
				ellipse.setStroke(currentColor);

				if (fillCheckBox.isSelected()) {
					ellipse.setFill(currentColor);
				} else {
					ellipse.setFill(null);
				}

				drawPane.getChildren().add(ellipse);
			}

		} else if (penStatus == PenStatus.ERASE) {

			int numberOfNodes = drawPane.getChildren().size();

			if (numberOfNodes > 0) {
				for (int i = 0; i < numberOfNodes; i++) {
					// System.out.println("numberOfNodes: " + numberOfNodes + " i:" + i);
					// System.out.println("node: " + drawPane.getChildren().get(i));
					boolean removeNode = drawPane.getChildren().get(i).contains(x, y);
					if (removeNode) {
						drawPane.getChildren().remove(i);
						numberOfNodes--;
					}
				}
			}
		}
	}

	private void handleMouseClicks(MouseEvent event) {

		isPenDown = !isPenDown;
		if (isPenDown) {
			penStatusHBoxText.setText("draw");
			penStatus = PenStatus.DRAW;
			currentColor = previousColor;
		} else {
			penStatusHBoxText.setText("off");
			penStatus = PenStatus.OFF;
			currentColor = Color.TRANSPARENT;
		}
	}

	public static void main(String[] args) {
		launch(args);
	}

}
