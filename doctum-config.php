<?php

use Doctum\Doctum;
use Symfony\Component\Finder\Finder;

$iterator = Finder::create()
    ->files()
    ->name("*.php")
    ->exclude(".changes")
    ->exclude("docker")
    ->exclude("runtime")
    ->exclude("tests")
    ->exclude("src")
    ->in("/work/repository/");

return new Doctum($iterator);
