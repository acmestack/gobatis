/*
 * Copyright (c) 2022, OpeningO
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package session

//
//import "github.com/xfali/gobatis/errors"
//
//type Propagation interface {
//    Begin() error
//    Commit() error
//    Rollback() error
//}
//
//type PropagationBase struct {
//    sess SqlSession
//}
//
//func (p *PropagationBase) SetSession(sess SqlSession) {
//    p.sess = sess
//}
//
////当前没有事务则新建事务，有则加入当前事务
//type PropagationRequired struct {
//    PropagationBase
//
//    beginCount int32
//}
//
//func (p *PropagationRequired) Begin() error {
//    if p.beginCount > 0 {
//        return nil
//    }
//
//    err := p.sess.Begin()
//    if err != nil {
//        return err
//    }
//
//    p.beginCount++
//    return nil
//}
//
//func (p *PropagationRequired) Commit() error {
//    if p.beginCount == 0 {
//        return errors.TRANSACTION_WITHOUT_BEGIN
//    }
//
//    p.beginCount--
//    if p.beginCount == 0 {
//        return p.sess.Commit()
//    }
//    return nil
//}
//
//func (p *PropagationRequired) Rollback() error {
//    if p.beginCount == 0 {
//        return errors.TRANSACTION_WITHOUT_BEGIN
//    }
//
//    p.beginCount--
//    if p.beginCount == 0 {
//        return p.sess.Rollback()
//    }
//    return nil
//}
//
////支持当前事务，如果当前没有事务则以非事务方式执行
//type PropagationSupports struct {
//    PropagationBase
//}
//
//func (p *PropagationSupports) Begin() error {
//
//}
//
//func (p *PropagationSupports) Commit() error {
//    if p.beginCount == 0 {
//        return errors.TRANSACTION_WITHOUT_BEGIN
//    }
//
//    p.beginCount--
//    if p.beginCount == 0 {
//        return p.sess.Commit()
//    }
//    return nil
//}
//
//func (p *PropagationSupports) Rollback() error {
//    if p.beginCount == 0 {
//        return errors.TRANSACTION_WITHOUT_BEGIN
//    }
//
//    p.beginCount--
//    if p.beginCount == 0 {
//        return p.sess.Rollback()
//    }
//    return nil
//}
//
////使用当前事务，如果没有则panic
//type PropagationMandatory struct {
//    PropagationBase
//}
//
////新建事务，如果当前有事务则把当前事务挂起
//type PropagationRequiredNew struct {
//    PropagationBase
//}
//
////以非事务方式执行操作，如果当前存在事务，就把当前事务挂起
//type PropagationNotSupported struct {
//    PropagationBase
//}
//
////以非事务的方式执行，如果当前有事务则panic
//type PropagationNever struct {
//    PropagationBase
//}
//
////如果当前存在事务，则在嵌套事务内执行。如果当前没有事务，则执行与PROPAGATION_REQUIRED类似的操作
//type PropagationNested struct {
//    PropagationBase
//}
