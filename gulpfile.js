#!/usr/bin/env gulp
var gulp = require('gulp');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var cleanCSS = require('gulp-clean-css');
var del = require('del');

// Assets path & destination definition
var paths = {
  scripts: {
    src: [
      'node_modules/jquery/dist/jquery.slim.js',
      'node_modules/popper.js/dist/umd/popper.js',
      'node_modules/bootstrap/dist/js/bootstrap.js',
      'node_modules/datatables/media/js/jquery.dataTables.js',
      'node_modules/datatables.net-bs4/js/dataTables.bootstrap4.js'
    ],
    dest: 'dist/js/'
  },
  styles: {
    src: [
      'node_modules/bootstrap/dist/css/bootstrap.css',
      'node_modules/datatables.net-bs4/css/dataTables.bootstrap4.css',
      'node_modules/open-iconic/font/css/open-iconic-bootstrap.css'
    ],
    dest: 'dist/css/'
  },
  fonts: {
    src: [
      'node_modules/open-iconic/font/fonts/*'
    ],
    dest: 'dist/fonts/'
  }
};

// Tasks definitions
function clean() {
  return del([paths.styles.dest, paths.scripts.dest, paths.fonts.dest]);
}

function styles() {
  return gulp.src(paths.styles.src)
  .pipe(concat('dist.css'))
  .pipe(cleanCSS())
  .pipe(gulp.dest(paths.styles.dest));
}

function scripts() {
  return gulp.src(paths.scripts.src)
    .pipe(uglify())
    .pipe(concat('dist.js'))
    .pipe(gulp.dest(paths.scripts.dest));
}

function fonts() {
  return gulp.src(paths.fonts.src)
    .pipe(gulp.dest(paths.fonts.dest));
}

var build = gulp.series(clean, gulp.parallel(styles, scripts, fonts));

// Tasks assignment
gulp.task('scripts', scripts);
gulp.task('styles', styles);
gulp.task('fonts', fonts);
gulp.task('clean', clean);
gulp.task('build', build);
gulp.task('default', build);
