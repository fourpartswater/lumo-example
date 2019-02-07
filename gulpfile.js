'use strict';

const gulp = require('gulp');
const babel = require('babelify');
const browserify = require('browserify');
const concat = require('gulp-concat');
const eslint = require('gulp-eslint');
const gulpgo = require('gulp-go');
const source = require('vinyl-source-stream');
const del = require('del');

const project = 'lumo-example';
const paths = {
	root: 'webapp/app.js',
	index: 'webapp/index.html',
	images: 'webapp/*.png',
	scripts: [ 'webapp/**/*.js' ],
	styles: [ 'webapp/**/*.css' ],
	go: [ '**/*.go' ],
	build: 'build'
};

function logError(err) {
	if (err instanceof SyntaxError) {
		console.error('Syntax Error:');
		console.error(err.message);
		console.error(err.codeFrame);
	} else {
		console.error(err.message);
	}
}

function handleError(err) {
	logError(err);
	this.emit('end');
}

function clean() {
	return del([ paths.build ]);
}

function lint() {
	return gulp.src(paths.scripts)
		.pipe(eslint())
		.pipe(eslint.format())
		.pipe(eslint.failAfterError());
}

function styles() {
	return gulp.src(paths.styles)
		.pipe(concat(`${project}.css`))
		.pipe(gulp.dest(paths.build));
}

function scripts() {
	return browserify(paths.root, {
		debug: true,
		standalone: project
	}).transform(babel, {
		presets: ['@babel/preset-env']
	})
	.bundle()
	.on('error', handleError)
	.pipe(source(`${project}.js`))
	.pipe(gulp.dest(paths.build));
}

function copyIndex() {
	return gulp.src(paths.index)
		.pipe(gulp.dest(paths.build));
}

function copyImages() {
	return gulp.src(paths.images)
		.pipe(gulp.dest(paths.build));
}

let go;
function serve() {
	go = gulpgo.run(paths.server, [], {
		stdio: 'inherit'
	});
}

function watch() {
	gulp.watch(paths.scripts, scripts);
	gulp.watch(paths.styles, styles);
	gulp.watch(paths.images, copyImages);
	gulp.watch(paths.index, copyIndex);
	gulp.watch(paths.go, () => {
		go.restart();
	});
}

const build = gulp.series(clean, gulp.parallel(scripts, styles, copyIndex, copyImages));

exports.clean = clean;
exports.lint = lint;
exports.build = build;
exports.serve = serve;

exports.default = gulp.series(build, gulp.parallel(watch, serve));
